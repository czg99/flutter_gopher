package bridge

/*
#include "include/bridge.h"
*/
import "C"
import (
	"unsafe"
)

func CallNativeMethod(method string, data []byte) []byte {
	if fgMethodHandle == nil {
		return nil
	}
	request := C.FgRequest{
		method: mapFgDataFromString(method),
		data:   mapFgDataFromBytes(data),
	}

	response := C.FgResponse{}
	C.call_fg_method_handle(fgMethodHandle, request, &response)
	defer freeFgResponse(&response)
	return mapFgDataToBytes(response.data)
}

func CallDartMethod(method string, data []byte) {
	request := C.FgRequest{
		method: mapFgDataFromString(method),
		data:   mapFgDataFromBytes(data),
	}
	fg_call_dart_method(request)
}

//export fg_empty_data
func fg_empty_data() C.FgData {
	return C.FgData{}
}

//export fg_empty_request
func fg_empty_request() C.FgRequest {
	return C.FgRequest{}
}

//export fg_empty_response
func fg_empty_response() C.FgResponse {
	return C.FgResponse{}
}

var global_port C.int64_t = 0

//export fg_init_dart_api
func fg_init_dart_api(api unsafe.Pointer, port C.int64_t) {
	init_dart_api(api)
	global_port = port
}

var fgMethodHandle C.FgMethodHandle = nil

//export fg_init_method_handle
func fg_init_method_handle(handle C.FgMethodHandle) {
	fgMethodHandle = handle
}

//export fg_call_dart_method
func fg_call_dart_method(request C.FgRequest) {
	if global_port == 0 {
		freeFgRequest(&request)
		return
	}
	sendToPort(global_port, unsafe.Pointer(cValueToPtr(request)))
}

//export fg_call_go_method
func fg_call_go_method(request C.FgRequest) C.FgResponse {
	method := mapFgDataToString(request.method)
	data := mapFgDataToBytes(request.data)
	freeFgRequest(&request)

	result := callGoMethod(method, data)
	response := C.FgResponse{
		data: mapFgDataFromBytes(result),
	}
	return response
}

//export fg_call_go_method_async
func fg_call_go_method_async(port C.int64_t, request C.FgRequest) {
	go func() {
		response := fg_call_go_method(request)
		sendToPort(port, unsafe.Pointer(cValueToPtr(response)))
	}()
}

//export fg_call_native_method
func fg_call_native_method(request C.FgRequest) C.FgResponse {
	if fgMethodHandle == nil {
		freeFgRequest(&request)
		return C.FgResponse{}
	}

	response := C.FgResponse{}
	C.call_fg_method_handle(fgMethodHandle, request, &response)
	return response
}

//export fg_call_native_method_async
func fg_call_native_method_async(port C.int64_t, request C.FgRequest) {
	go func() {
		response := fg_call_native_method(request)
		sendToPort(port, unsafe.Pointer(cValueToPtr(response)))
	}()
}

//export enforce_binding
func enforce_binding() {
	var ptr uintptr
	ptr ^= uintptr(unsafe.Pointer(C.fg_empty_data))
	ptr ^= uintptr(unsafe.Pointer(C.fg_empty_request))
	ptr ^= uintptr(unsafe.Pointer(C.fg_empty_response))
	ptr ^= uintptr(unsafe.Pointer(C.fg_init_dart_api))
	ptr ^= uintptr(unsafe.Pointer(C.fg_init_method_handle))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_dart_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_go_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_go_method_async))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_native_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_native_method_async))
}

func mapFgDataFromString(from string) C.FgData {
	return mapFgDataFromBytes([]byte(from))
}

func mapFgDataToString(from C.FgData) string {
	data := mapFgDataToBytes(from)
	if data == nil {
		return ""
	}
	return string(data)
}

func mapFgDataFromBytes(from []byte) C.FgData {
	if from == nil {
		return C.FgData{}
	}
	data := C.CBytes(from)
	size := C.int(len(from))
	return C.FgData{
		data: data,
		size: size,
	}
}

func mapFgDataToBytes(from C.FgData) []byte {
	if from.data == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(from.data), C.int(from.size))
}

func cValueToPtr[T any](value T) *T {
	size := unsafe.Sizeof(value)
	data := C.malloc(C.size_t(size))
	*(*T)(data) = value
	return (*T)(data)
}

func freeFgData(value *C.FgData) {
	if value.data != nil {
		C.free(value.data)
		value.data = nil
		value.size = 0
	}
}

func freeFgRequest(request *C.FgRequest) {
	freeFgData(&request.method)
	freeFgData(&request.data)
}

func freeFgResponse(response *C.FgResponse) {
	freeFgData(&response.data)
}
