#include <jni.h>
#include "../include/bridge.h"

#ifndef FG_BRIDGE_JNI_H
#define FG_BRIDGE_JNI_H

JNIEXPORT void JNICALL Java_{{.JniPackagePath}}_FgBridge_nativeInit(JNIEnv* env, jclass clazz);
JNIEXPORT jobject JNICALL Java_{{.JniPackagePath}}_FgBridge_nativeCallGoMethod(JNIEnv* env, jclass clazz, jint method, jbyteArray data);
JNIEXPORT void JNICALL Java_{{.JniPackagePath}}_FgBridge_nativeCallDartMethod(JNIEnv* env, jclass clazz, jint method, jbyteArray data);

#endif
