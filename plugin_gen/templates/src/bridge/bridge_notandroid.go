//go:build !android
// +build !android

package bridge

/*
#include "bridge.h"
*/
import "C"
import "unsafe"

func CallMethod(method string, data []byte) []byte {
	if fgMethodHandle == nil {
		return nil
	}
	c_param := C.FgPacket{
		method:     C.CString(method),
		method_len: C.int(len(method)),
	}

	if data != nil {
		c_param.data = unsafe.Pointer(C.CBytes(data))
		c_param.data_len = C.int(len(data))
	}

	c_result := C.call_fg_method_handle(fgMethodHandle, c_param)
	defer C.free(unsafe.Pointer(c_result.method))
	defer C.free(unsafe.Pointer(c_result.data))
	if c_result.data == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(c_result.data), C.int(c_result.data_len))
}

var fgMethodHandle C.FgMethodHandle = nil

//export fg_init_method_handle
func fg_init_method_handle(handle C.FgMethodHandle) {
	fgMethodHandle = handle
}
