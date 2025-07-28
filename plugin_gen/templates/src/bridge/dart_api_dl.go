package bridge

/*
#include <stdlib.h>
#include <stdint.h>
#include "include/dart_api_dl.h"
#include "include/dart_api_dl.c"

bool GoDart_PostCObject(Dart_Port send_port, void* data) {
  Dart_CObject dart_object;
  dart_object.type = Dart_CObject_kInt64;
  dart_object.value.as_int64 = (int64_t)data;
  return Dart_PostCObject_DL(send_port, &dart_object);
}
*/
import "C"
import (
	"unsafe"
)

func init_dart_api(api unsafe.Pointer) {
	if C.Dart_InitializeApiDL(api) != 0 {
		panic("failed to initialize Dart DL C API: version mismatch. " +
			"must update include/ to match Dart SDK version")
	}
}

func sendToPort(port C.int64_t, data unsafe.Pointer) bool {
	return C.GoDart_PostCObject(C.Dart_Port_DL(port), data) == true
}
