#include <napi/native_api.h>
#include "../include/bridge.h"

static napi_ref g_platformCallbackRef = NULL;
static napi_env g_env = NULL;

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
        data.size = (int)length;
        memcpy(data.data, buffer, length);
    }
    return data;
}

/*************** FgData -> JS ***************/
static napi_value fgdata_to_js(napi_env env, FgData data) {
    napi_value arraybuffer = NULL;
    void *buffer = NULL;

    if (data.size <= 0) {
        napi_create_arraybuffer(env, 0, &buffer, &arraybuffer);
        return arraybuffer;
    }

    napi_create_arraybuffer(env, data.size, &buffer, &arraybuffer);
    memcpy(buffer, data.data, data.size);
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

static napi_value napi_init_platform_method_handle(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value args[1] = {NULL};

    napi_get_cb_info(env, info, &argc, args, NULL, NULL);
    if (argc < 1) {
        return NULL;
    }

    napi_valuetype valueType;
    napi_typeof(env, args[0], &valueType);
    if (valueType != napi_function) {
        return NULL;
    }

    if (g_platformCallbackRef != NULL) {
        napi_delete_reference(env, g_platformCallbackRef);
        g_platformCallbackRef = NULL;
    }

    napi_create_reference(env, args[0], 1, &g_platformCallbackRef);

    g_env = env;
    return NULL;
}

static napi_value napi_call_go_method(napi_env env, napi_callback_info info) {
    size_t argc = 1;
    napi_value args[1] = {NULL};
    napi_get_cb_info(env, info, &argc, args, NULL, NULL);

    FgRequest req = js_to_fgrequest(env, args[0]);
    FgResponse resp = fg_call_go_method_{{.ID}}(req);

    napi_value result = fgresponse_to_js(env, resp);

    free_fgdata(&resp.data);
    free_fgdata(&resp.error);
    return result;
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

DLLEXPORT void napi_onLoad()
{    
    napi_module_register(&{{.LibName}}Module);
}

// __attribute__((constructor)) void RegisterModule(void)
// {
//     napi_module_register(&{{.LibName}}Module);
// }

static FgResponse fg_call_platform_method(FgRequest request) {
    FgResponse resp = {};
    if (g_platformCallbackRef == nullptr || g_env == nullptr) {
        return resp;
    }

    napi_handle_scope scope = NULL;
    napi_open_handle_scope(g_env, &scope);

    napi_value callback = NULL;
    napi_get_reference_value(g_env, g_platformCallbackRef, &callback);

    napi_value global = NULL;
    napi_get_global(g_env, &global);

    napi_value jsRequest = NULL;
    napi_create_object(g_env, &jsRequest);

    napi_value method = NULL;
    napi_create_int32(g_env, request.method, &method);
    napi_set_named_property(g_env, jsRequest, "method", method);

    napi_value data = fgdata_to_js(g_env, request.data);
    napi_set_named_property(g_env, jsRequest, "data", data);

    napi_value result = NULL;
    napi_call_function(g_env, global, callback, 1, &jsRequest, &result);

    if (result != nullptr) {
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
