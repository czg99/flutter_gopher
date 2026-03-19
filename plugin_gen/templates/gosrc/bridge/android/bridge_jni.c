#include "bridge_jni.h"

static JavaVM* g_vm = NULL;
static jclass g_bridgeClass = NULL;
static jmethodID g_handleMid = NULL;

//jbyteArray <-> FgData
static FgData jni_to_fgdata(JNIEnv* env, jbyteArray array) {
    if (!array) return (FgData){NULL, 0};
    int32_t len = (*env)->GetArrayLength(env, array);
    if (len == 0) return (FgData){NULL, 0};
    void* buf = malloc(len);
    (*env)->GetByteArrayRegion(env, array, 0, len, (jbyte*)buf);
    return (FgData){buf, len};
}

static jbyteArray fgdata_to_jni(JNIEnv* env, FgData data) {
    if (!data.data || data.size <= 0) return NULL;
    jbyteArray array = (*env)->NewByteArray(env, data.size);
    (*env)->SetByteArrayRegion(env, array, 0, data.size, (jbyte*)data.data);
    free(data.data);
    return array;
}

static void bridge_callback(FgRequest req, FgResponse* resp) {
    JNIEnv* env = NULL;
    if ((*g_vm)->GetEnv(g_vm, (void**)&env, JNI_VERSION_1_6) != JNI_OK) {
        if ((*g_vm)->AttachCurrentThread(g_vm, &env, NULL) != JNI_OK) return;
    }

    jbyteArray jData = fgdata_to_jni(env, req.data);
    
    jobject jRespObj = (*env)->CallStaticObjectMethod(env, g_bridgeClass, g_handleMid, req.method, jData);
    
    if (jRespObj) {
        jclass resClazz = (*env)->GetObjectClass(env, jRespObj);
        jfieldID dataFid = (*env)->GetFieldID(env, resClazz, "data", "[B");
        jfieldID errorFid = (*env)->GetFieldID(env, resClazz, "error", "[B");
        
        jbyteArray resData = (jbyteArray)(*env)->GetObjectField(env, jRespObj, dataFid);
        jbyteArray resError = (jbyteArray)(*env)->GetObjectField(env, jRespObj, errorFid);
        
        if (resData) {
            resp->data = jni_to_fgdata(env, resData);
        }
        if (resError) {
            resp->error = jni_to_fgdata(env, resError);
        }
    }
}

JNIEXPORT void JNICALL Java_{{.JniPackagePath}}_FgBridge_nativeInit(JNIEnv* env, jclass clazz) {
    (*env)->GetJavaVM(env, &g_vm);
    g_bridgeClass = (jclass)(*env)->NewGlobalRef(env, clazz);
    g_handleMid = (*env)->GetStaticMethodID(env, clazz, "onNativeHandle", "(I[B)L{{.JavaPackagePath}}/FgJniResponse;");
    
    fg_init_platform_method_handle_{{.ID}}(bridge_callback);
}

JNIEXPORT jobject JNICALL Java_{{.JniPackagePath}}_FgBridge_nativeCallGoMethod(JNIEnv* env, jclass clazz, jint method, jbyteArray data) {
    FgRequest req = { (int32_t)method, jni_to_fgdata(env, data) };
    FgResponse res = fg_call_go_method_{{.ID}}(req);

    jbyteArray jData = fgdata_to_jni(env, res.data);
    jbyteArray jError = fgdata_to_jni(env, res.error);

    jclass resClazz = (*env)->FindClass(env, "{{.JavaPackagePath}}/FgJniResponse");
    jmethodID init = (*env)->GetMethodID(env, resClazz, "<init>", "([B[B)V");
    return (*env)->NewObject(env, resClazz, init, jData, jError);
}

JNIEXPORT void JNICALL Java_{{.JniPackagePath}}_FgBridge_nativeCallDartMethod(JNIEnv* env, jclass clazz, jint method, jbyteArray data) {
    FgRequest req = { (int32_t)method, jni_to_fgdata(env, data) };
    fg_call_dart_method_{{.ID}}(req);
}
