#include <napi/native_api.h>
#include <pthread.h>
#include <future>
#include "hilog/log.h"

extern "C" {
    #include "../include/bridge.h"
}

static void callJsCallback(napi_env env, napi_value js_cb, void* context, void* data);

static napi_threadsafe_function g_tsfn = nullptr;
static napi_ref g_platformCallbackRef = nullptr;
static napi_env g_env = nullptr;
static pthread_t g_main_thread_id = 0;

typedef struct {
    napi_async_work asyncWork;
    napi_deferred deferred;
    napi_ref callback;
    FgRequest request;
    FgResponse response;
} FgAsyncData;

typedef struct {
    FgRequest request;
    std::promise<FgResponse*> promise;
} FgThreadObject;

/*************** JS <-> FgData ***************/
static FgData js_to_fgdata(napi_env env, napi_value jsData) {
    FgData data = {};

    bool isArrayBuffer = false;
    napi_is_arraybuffer(env, jsData, &isArrayBuffer);
    if (!isArrayBuffer)
        return data;

    void* buffer = nullptr;
    size_t length = 0;
    napi_get_arraybuffer_info(env, jsData, &buffer, &length);

    if (length > 0) {
        data.data = malloc(length);
        data.size = (int32_t)length;
        memcpy(data.data, buffer, length);
    }
    return data;
}

static napi_value fgdata_to_js(napi_env env, FgData data) {
    napi_value arraybuffer = nullptr;
    void* buffer = nullptr;

    if (data.data == nullptr) {
        napi_create_arraybuffer(env, 0, &buffer, &arraybuffer);
        return arraybuffer;
    }

    napi_create_arraybuffer(env, data.size, &buffer, &arraybuffer);
    memcpy(buffer, data.data, data.size);
    free_fgdata(&data);
    return arraybuffer;
}

/*************** FgRequest <-> JS ***************/
static napi_value fgrequest_to_js(napi_env env, FgRequest request) {
    napi_value method = nullptr;
    napi_create_int32(env, request.method, &method);
    napi_value data = fgdata_to_js(env, request.data);

    napi_value jsRequest = nullptr;
    napi_create_object(env, &jsRequest);
    napi_set_named_property(env, jsRequest, "method", method);
    napi_set_named_property(env, jsRequest, "data", data);
    return jsRequest;
}

static FgRequest js_to_fgrequest(napi_env env, napi_value jsRequest) {
    napi_value methodValue = nullptr;
    napi_value dataValue = nullptr;
    napi_get_named_property(env, jsRequest, "method", &methodValue);
    napi_get_named_property(env, jsRequest, "data", &dataValue);

    FgRequest request = {};
    napi_get_value_int32(env, methodValue, &request.method);
    request.data = js_to_fgdata(env, dataValue);
    return request;
}

/*************** FgResponse <-> JS ***************/
static napi_value fgresponse_to_js(napi_env env, FgResponse response) {
    napi_value data = fgdata_to_js(env, response.data);
    napi_value error = fgdata_to_js(env, response.error);

    napi_value jsResponse = nullptr;
    napi_create_object(env, &jsResponse);
    napi_set_named_property(env, jsResponse, "data", data);
    napi_set_named_property(env, jsResponse, "error", error);
    return jsResponse;
}

static FgResponse js_to_fgresponse(napi_env env, napi_value jsResponse) {
    napi_value jsData = nullptr;
    napi_value jsError = nullptr;
    napi_get_named_property(env, jsResponse, "data", &jsData);
    napi_get_named_property(env, jsResponse, "error", &jsError);
    
    FgResponse response = {};
    response.data = js_to_fgdata(env, jsData);
    response.error = js_to_fgdata(env, jsError);
    return response;
}


static int is_main_thread(void) {
    return pthread_equal(pthread_self(), g_main_thread_id);
}

static napi_value napi_init_platform_method_handle(napi_env env, napi_callback_info info) {
    g_main_thread_id = pthread_self();
    g_env = env;

    size_t argc = 1;
    napi_value args[1] = {nullptr};

    napi_get_cb_info(env, info, &argc, args, nullptr, nullptr);
    if (argc < 1) {
        return nullptr;
    }

    napi_value callback = args[0];

    napi_valuetype valueType = napi_undefined;
    napi_typeof(env, callback, &valueType);
    if (valueType != napi_function) {
        return nullptr;
    }

    if (g_tsfn != nullptr) {
        napi_release_threadsafe_function(g_tsfn, napi_tsfn_abort);
        g_tsfn = nullptr;
    }

    napi_value resourceName = nullptr;
    napi_create_string_utf8(env, "{{.LibName}}_napi_init_platform_method_handle", NAPI_AUTO_LENGTH, &resourceName);

    napi_create_threadsafe_function(
        env,
        callback,
        nullptr,
        resourceName,
        8,
        1,
        nullptr,
        nullptr,
        nullptr,
        callJsCallback,
        &g_tsfn
    );

    if (g_platformCallbackRef != nullptr) {
        napi_delete_reference(env, g_platformCallbackRef);
        g_platformCallbackRef = nullptr;
    }

    napi_create_reference(env, callback, 1, &g_platformCallbackRef);
    return nullptr;
}

static napi_value napi_call_go_method(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value args[1] = {nullptr};
    napi_get_cb_info(env, info, &argc, args, nullptr, nullptr);

    FgRequest req = js_to_fgrequest(env, args[0]);
    FgResponse resp = fg_call_go_method_{{.ID}}(req);

    return fgresponse_to_js(env, resp);
}

static void napi_call_go_method_async_execute(napi_env env, void *data) {
    FgAsyncData* adata = (FgAsyncData*)data;
    adata->response = fg_call_go_method_{{.ID}}(adata->request);;
}

static void napi_call_go_method_async_complete(napi_env env, napi_status status, void *data) {
    FgAsyncData* adata = (FgAsyncData*)data;
    napi_value response = fgresponse_to_js(env, adata->response);
    napi_resolve_deferred(env, adata->deferred, response);
    napi_delete_async_work(env, adata->asyncWork);
    free(adata);
}

static napi_value napi_call_go_method_async(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value args[1] = {nullptr};
    napi_get_cb_info(env, info, &argc, args, nullptr, nullptr);

    napi_value promise = nullptr;
    napi_deferred deferred = nullptr;
    napi_create_promise(env, &deferred, &promise);

    FgAsyncData* adata = (FgAsyncData*)calloc(1, sizeof(FgAsyncData));
    adata->request = js_to_fgrequest(env, args[0]);
    adata->deferred = deferred;

    napi_value resourceName = nullptr;
    napi_create_string_utf8(env, "{{.LibName}}_napi_call_go_method_async", NAPI_AUTO_LENGTH, &resourceName);

    napi_create_async_work(env, nullptr, resourceName, napi_call_go_method_async_execute, napi_call_go_method_async_complete, adata, &adata->asyncWork);
    napi_queue_async_work(env, adata->asyncWork);
    return promise;
}


static napi_value napi_call_dart_method(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value args[1] = {nullptr};
    napi_get_cb_info(env, info, &argc, args, nullptr, nullptr);

    FgRequest req = js_to_fgrequest(env, args[0]);
    fg_call_dart_method_{{.ID}}(req);
    return nullptr;
}

static napi_value Init(napi_env env, napi_value exports) {
    napi_property_descriptor desc[] = {
        {"initPlatformMethodHandle", 0, napi_init_platform_method_handle, 0, 0, 0, napi_default, 0},
        {"callGoMethod", 0, napi_call_go_method, 0, 0, 0, napi_default, 0},
        {"callGoMethodAsync", 0, napi_call_go_method_async, 0, 0, 0, napi_default, 0},
        {"callDartMethod", 0, napi_call_dart_method, 0, 0, 0, napi_default, 0},
    };
    napi_define_properties(env, exports, sizeof(desc) / sizeof(desc[0]), desc);
    return exports;
}

static napi_module {{.LibName}}Module = {
    .nm_version = 1,
    .nm_flags = 0,
    .nm_filename = nullptr,
    .nm_register_func = Init,
    .nm_modname = "{{.LibName}}",
    .nm_priv = ((void*)0),
    .reserved = { 0 },
};

__attribute__((constructor)) static void RegisterModule(void) {
    napi_module_register(&{{.LibName}}Module);
}

static napi_value resolved_callback(napi_env env, napi_callback_info info) {
    void *data = nullptr;
    size_t argc = 1;
    napi_value argv[1];
    if (napi_get_cb_info(env, info, &argc, argv, nullptr, &data) != napi_ok) {
        return nullptr;
    }
    FgResponse response = js_to_fgresponse(env, argv[0]);
    FgResponse* responsePtr = (FgResponse*)malloc(sizeof(FgResponse));
    *responsePtr = response;
    reinterpret_cast<std::promise<FgResponse*>*>(data)->set_value(responsePtr);
    return nullptr;
}

static napi_value rejected_callback(napi_env env, napi_callback_info info) {
    void *data = nullptr;
    if (napi_get_cb_info(env, info, nullptr, nullptr, nullptr, &data) != napi_ok) {
        return nullptr;
    }
    reinterpret_cast<std::promise<FgResponse*>*>(data)->set_exception(
        std::make_exception_ptr(std::runtime_error("Error in jsCallback")));
    return nullptr;
}

static void fg_call_platform_method_promise(napi_env env, napi_value jsCallback, FgRequest request, std::promise<FgResponse*>* promise) {
    napi_value jsRequest = fgrequest_to_js(env, request);

    napi_value jsPromise = nullptr;
    napi_call_function(env, nullptr, jsCallback, 1, &jsRequest, &jsPromise);

    napi_value thenFunc = nullptr;
    if (napi_get_named_property(env, jsPromise, "then", &thenFunc) != napi_ok) {
        promise->set_exception(std::make_exception_ptr(std::runtime_error("Error in jsCallback")));
        return;
    }

    napi_value resolvedCallback;
    napi_value rejectedCallback;
    napi_create_function(env, "{{.LibName}}_resolved_callback", NAPI_AUTO_LENGTH, resolved_callback, promise, &resolvedCallback);
    napi_create_function(env, "{{.LibName}}_rejected_callback", NAPI_AUTO_LENGTH, rejected_callback, promise, &rejectedCallback);
    napi_value argv[2] = {resolvedCallback, rejectedCallback};
    napi_call_function(env, jsPromise, thenFunc, 2, argv, nullptr);
}

extern "C" {
    FgResponse fg_call_platform_method(FgRequest request) {
        FgResponse response = {};

        if (g_platformCallbackRef == nullptr || g_env == nullptr) {
            return response;
        }

        napi_value jsCallback;
        napi_get_reference_value(g_env, g_platformCallbackRef, &jsCallback);

        std::promise<FgResponse*> promise;
        auto future = promise.get_future();

        fg_call_platform_method_promise(g_env, jsCallback, request, &promise);

        try {
            auto responsePtr = future.get();
            response = *responsePtr;
            free(responsePtr);
        } catch (const std::exception &e) {
        }
        return response;
    }

    FgResponse fg_call_platform_method_safe(FgRequest request) {
        if (is_main_thread()) {
            return fg_call_platform_method(request);
        }

        if (g_tsfn == nullptr) {
            return (FgResponse){};
        }
    
        FgThreadObject obj = {};
        obj.request = request;

        auto future = obj.promise.get_future();

        napi_call_threadsafe_function(g_tsfn, &obj, napi_tsfn_blocking);

        FgResponse response = {};
        try {
            auto responsePtr = future.get();
            response = *responsePtr;
            free(responsePtr);
        } catch (const std::exception &e) {
        }
        return response;
    }
}

static void callJsCallback(napi_env env, napi_value js_cb, void* context, void* data) {
    FgThreadObject *obj = (FgThreadObject *)data;
    fg_call_platform_method_promise(env, js_cb, obj->request, &obj->promise);
}
