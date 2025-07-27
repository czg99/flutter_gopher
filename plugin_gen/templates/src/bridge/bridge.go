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
	packet := C.FgPacket{
		method: mapFgDataFromString(method),
		data:   mapFgDataFromBytes(data),
	}

	c_result := C.FgPacket{}
	C.call_fg_method_handle(fgMethodHandle, packet, &c_result)
	defer freeFgData(&c_result.method)
	defer freeFgData(&c_result.data)
	return mapFgDataToBytes(c_result.data)
}

//export fg_init_dart_api
func fg_init_dart_api(api unsafe.Pointer) {
	init_dart_api(api)
}

//export fg_empty_packet
func fg_empty_packet() C.FgPacket {
	return C.FgPacket{}
}

//export fg_empty_data
func fg_empty_data() C.FgData {
	return C.FgData{}
}

//export fg_call_go_method
func fg_call_go_method(packet C.FgPacket) C.FgPacket {
	method := mapFgDataToString(packet.method)
	data := mapFgDataToBytes(packet.data)
	freeFgData(&packet.data)

	result := callGoMethod(method, data)
	packet.data = mapFgDataFromBytes(result)
	return packet
}

//export fg_call_go_method_async
func fg_call_go_method_async(port int64, packet C.FgPacket) {
	go func() {
		result := fg_call_go_method(packet)
		size := unsafe.Sizeof(result)
		data := C.malloc(C.size_t(size))
		*(*C.FgPacket)(data) = result
		sendToPort(port, data)
	}()
}

//export fg_call_native_method
func fg_call_native_method(packet C.FgPacket) C.FgPacket {
	if fgMethodHandle == nil {
		freeFgData(&packet.data)
		return packet
	}

	c_result := C.FgPacket{}
	C.call_fg_method_handle(fgMethodHandle, packet, &c_result)
	return c_result
}

//export fg_call_native_method_async
func fg_call_native_method_async(port int64, packet C.FgPacket) {
	go func() {
		result := fg_call_native_method(packet)
		size := unsafe.Sizeof(result)
		data := C.malloc(C.size_t(size))
		*(*C.FgPacket)(data) = result
		sendToPort(port, data)
	}()
}

var fgMethodHandle C.FgMethodHandle = nil

//export fg_init_method_handle
func fg_init_method_handle(handle C.FgMethodHandle) {
	fgMethodHandle = handle
}

//export enforce_binding
func enforce_binding() {
	var ptr uintptr
	ptr ^= uintptr(unsafe.Pointer(C.fg_init_dart_api))
	ptr ^= uintptr(unsafe.Pointer(C.fg_empty_data))
	ptr ^= uintptr(unsafe.Pointer(C.fg_empty_packet))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_go_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_go_method_async))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_native_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_native_method_async))
	ptr ^= uintptr(unsafe.Pointer(C.fg_init_method_handle))
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

func freeFgData(value *C.FgData) {
	if value.data != nil {
		C.free(value.data)
		value.data = nil
		value.size = 0
	}
}
