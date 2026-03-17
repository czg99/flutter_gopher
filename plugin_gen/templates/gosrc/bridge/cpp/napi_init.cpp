#include <napi/native_api.h>
#include <pthread.h>
#include "../include/bridge.h"

static void callJsCallback(napi_env env, napi_value js_cb, void* context, void* data);

static napi_threadsafe_function g_tsfn = NULL;
static napi_ref g_platformCallbackRef = NULL;
static napi_env g_env = NULL;
static pthread_t g_main_thread_id = 0;

static void free_fgdata(FgData *data) {
    if (data->data != NULL) {
        free(data->data);
        data->data = NULL;
    }
    data->size = 0;
}

/*************** JS -> FgData ***************/
static FgData js_to_fgdata(napi_env env, napi_value jsData) {
    FgData data = {};

    bool isArrayBuffer = false;
    napi_is_arraybuffer(env, jsData, &isArrayBuffer);
    if (!isArrayBuffer)
        return data;

    void *buffer = NULL;
    size_t length = 0;
    napi_get_arraybuffer_info(env, jsData, &buffer, &length);

    if (length > 0) {
        data.data = malloc(length);
        data.size = (int32_t)length;
        memcpy(data.data, buffer, length);
    }
    return data;
}

/*************** FgData -> JS ***************/
static napi_value fgdata_to_js(napi_env env, FgData data) {
    napi_value arraybuffer = NULL;
    void *buffer = NULL;

    if (data.data == NULL) {
        napi_create_arraybuffer(env, 0, &buffer, &arraybuffer);
        return arraybuffer;
    }

    napi_create_arraybuffer(env, data.size, &buffer, &arraybuffer);
    memcpy(buffer, data.data, data.size);
    free_fgdata(&data);
    return arraybuffer;
}

/*************** JS -> FgRequest ***************/
static FgRequest js_to_fgrequest(napi_env env, napi_value jsRequest) {
    FgRequest req = {};

    napi_value methodValue = NULL;
    napi_get_named_property(env, jsRequest, "method", &methodValue);
    napi_get_value_int32(env, methodValue, &req.method);

    napi_value dataValue = NULL;
    napi_get_named_property(env, jsRequest, "data", &dataValue);
    req.data = js_to_fgdata(env, dataValue);
    return req;
}

/*************** FgResponse -> JS ***************/
static napi_value fgresponse_to_js(napi_env env, FgResponse resp) {
    napi_value result = NULL;
    napi_create_object(env, &result);

    napi_value data = fgdata_to_js(env, resp.data);
    napi_value error = fgdata_to_js(env, resp.error);

    napi_set_named_property(env, result, "data", data);
    napi_set_named_property(env, result, "error", error);
    return result;
}

static int is_main_thread(void) {
    return pthread_equal(pthread_self(), g_main_thread_id);
}

static napi_value napi_init_platform_method_handle(napi_env env, napi_callback_info info) {
    g_main_thread_id = pthread_self();

    g_env = env;

    size_t argc = 1;
    napi_value args[1] = {NULL};

    napi_get_cb_info(env, info, &argc, args, NULL, NULL);
    if (argc < 1) {
        return NULL;
    }

    napi_value callback = args[0];

    napi_valuetype valueType = napi_undefined;
    napi_typeof(env, callback, &valueType);
    if (valueType != napi_function) {
        return NULL;
    }

    if (g_tsfn != NULL) {
        napi_release_threadsafe_function(g_tsfn, napi_tsfn_abort);
        g_tsfn = NULL;
    }

    napi_value async_resource_name = NULL;
    napi_create_string_utf8(env, "{{.LibName}}", NAPI_AUTO_LENGTH, &async_resource_name);

    napi_create_threadsafe_function(
        env,
        callback,
        NULL,
        async_resource_name,
        8,
        1,
        NULL,
        NULL,
        NULL,
        callJsCallback,
        &g_tsfn
    );

    if (g_platformCallbackRef != NULL) {
        napi_delete_reference(env, g_platformCallbackRef);
        g_platformCallbackRef = NULL;
    }

    napi_create_reference(env, callback, 1, &g_platformCallbackRef);
    return NULL;
}

static napi_value napi_call_go_method(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value args[1] = {NULL};
    napi_get_cb_info(env, info, &argc, args, NULL, NULL);

    FgRequest req = js_to_fgrequest(env, args[0]);
    FgResponse resp = fg_call_go_method_{{.ID}}(req);

    return fgresponse_to_js(env, resp);
}

static napi_value napi_call_dart_method(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value args[1] = {NULL};
    napi_get_cb_info(env, info, &argc, args, NULL, NULL);

    FgRequest req = js_to_fgrequest(env, args[0]);
    fg_call_dart_method_{{.ID}}(req);
    return NULL;
}

static napi_value Init(napi_env env, napi_value exports)
{
    napi_property_descriptor desc[] = {
        {"initPlatformMethodHandle", 0, napi_init_platform_method_handle, 0, 0, 0, napi_default, 0},
        {"callGoMethod", 0, napi_call_go_method, 0, 0, 0, napi_default, 0},
        {"callDartMethod", 0, napi_call_dart_method, 0, 0, 0, napi_default, 0},
    };
    napi_define_properties(env, exports, sizeof(desc) / sizeof(desc[0]), desc);
    return exports;
}

static napi_module {{.LibName}}Module = {
    .nm_version = 1,
    .nm_flags = 0,
    .nm_filename = NULL,
    .nm_register_func = Init,
    .nm_modname = "{{.LibName}}",
    .nm_priv = ((void*)0),
    .reserved = { 0 },
};

__attribute__((constructor)) static void RegisterModule(void)
{
    napi_module_register(&{{.LibName}}Module);
}

typedef struct {
    FgRequest req;
    FgResponse resp;
    bool done;
    pthread_mutex_t mutex;
    pthread_cond_t cond;
} FgThreadObject;

__attribute__((noinline)) static FgResponse fg_call_platform_method(FgRequest request) {
    FgResponse resp = {};
    if (g_platformCallbackRef == NULL || g_env == NULL) {
        return resp;
    }

    napi_handle_scope scope = NULL;
    napi_open_handle_scope(g_env, &scope);

    napi_value callback = NULL;
    napi_get_reference_value(g_env, g_platformCallbackRef, &callback);

    napi_value jsRequest = NULL;
    napi_create_object(g_env, &jsRequest);

    napi_value method = NULL;
    napi_create_int32(g_env, request.method, &method);
    napi_set_named_property(g_env, jsRequest, "method", method);

    napi_value data = fgdata_to_js(g_env, request.data);
    napi_set_named_property(g_env, jsRequest, "data", data);

    napi_value result = NULL;
    napi_call_function(g_env, NULL, callback, 1, &jsRequest, &result);

    if (result != NULL) {
        napi_value jsData = NULL;
        napi_value jsError = NULL;

        napi_get_named_property(g_env, result, "data", &jsData);
        napi_get_named_property(g_env, result, "error", &jsError);

        resp.data = js_to_fgdata(g_env, jsData);
        resp.error = js_to_fgdata(g_env, jsError);
    }

    napi_close_handle_scope(g_env, scope);
    return resp;
}

__attribute__((noinline)) static FgResponse fg_call_platform_method_safe(FgRequest request) {
    if (is_main_thread()) {
        return fg_call_platform_method(request);
    }

    if (g_tsfn == NULL) {
        return (FgResponse){};
    }

    FgThreadObject obj = {};
    obj.req = request;
    obj.resp = (FgResponse){};
    obj.done = false;
    pthread_mutex_init(&obj.mutex, NULL);
    pthread_cond_init(&obj.cond, NULL);

    napi_call_threadsafe_function(g_tsfn, &obj, napi_tsfn_blocking);

    pthread_mutex_lock(&obj.mutex);
    while (!obj.done) {
        pthread_cond_wait(&obj.cond, &obj.mutex);
    }
    pthread_mutex_unlock(&obj.mutex);
    return obj.resp;
}

static void callJsCallback(napi_env env, napi_value js_cb, void* context, void* data) {
    FgThreadObject *obj = (FgThreadObject *)data;
    obj->resp = fg_call_platform_method(obj->req);

    pthread_mutex_lock(&obj->mutex);
    obj->done = true;
    pthread_cond_signal(&obj->cond);
    pthread_mutex_unlock(&obj->mutex);
}
