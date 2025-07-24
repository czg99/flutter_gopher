package bridge

/*
#include "bridge.h"
*/
import "C"
import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var packetChan = make(chan C.FgPacket)

const minPortId int64 = 0xFF
const maxPortId int64 = 0xFFFFFFFFFFFF

var portMutex sync.Mutex
var nextPortId int64 = minPortId - 1

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
	defer freeFgData(c_result.method)
	defer freeFgData(c_result.data)
	return mapFgDataToBytes(c_result.data)
}

//export fg_next_port_id
func fg_next_port_id() C.int64_t {
	next := atomic.AddInt64(&nextPortId, 1)

	if next > maxPortId {
		portMutex.Lock()
		defer portMutex.Unlock()

		current := atomic.LoadInt64(&nextPortId)
		if current > maxPortId {
			atomic.StoreInt64(&nextPortId, minPortId)
			return C.int64_t(minPortId)
		} else {
			return C.int64_t(atomic.AddInt64(&nextPortId, 1))
		}
	}
	return C.int64_t(next)
}

//export fg_empty_packet
func fg_empty_packet() C.FgPacket {
	return C.FgPacket{}
}

//export fg_empty_data
func fg_empty_data() C.FgData {
	return C.FgData{}
}

//export fg_packet_loop
func fg_packet_loop() C.FgPacket {
	select {
	case result := <-packetChan:
		return result
	case <-time.After(time.Second):
		return fg_empty_packet()
	}
}

//export fg_call_go_method
func fg_call_go_method(packet C.FgPacket) C.FgPacket {
	defer freeFgData(packet.data)
	method, data := unpackFgPacket(packet)

	result := callGoMethod(method, data)
	return C.FgPacket{
		id:     packet.id,
		method: packet.method,
		data:   mapFgDataFromBytes(result),
	}
}

//export fg_call_go_method_async
func fg_call_go_method_async(packet C.FgPacket) {
	go func() {
		result := fg_call_go_method(packet)
		packetChan <- result
	}()
}

//export fg_call_native_method
func fg_call_native_method(packet C.FgPacket) C.FgPacket {
	if fgMethodHandle == nil {
		freeFgData(packet.data)
		return C.FgPacket{
			id:     packet.id,
			method: packet.method,
		}
	}

	c_result := C.FgPacket{id: packet.id}
	C.call_fg_method_handle(fgMethodHandle, packet, &c_result)
	return c_result
}

//export fg_call_native_method_async
func fg_call_native_method_async(packet C.FgPacket) {
	go func() {
		result := fg_call_native_method(packet)
		packetChan <- result
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
	ptr ^= uintptr(unsafe.Pointer(C.fg_empty_data))
	ptr ^= uintptr(unsafe.Pointer(C.fg_empty_packet))
	ptr ^= uintptr(unsafe.Pointer(C.fg_packet_loop))
	ptr ^= uintptr(unsafe.Pointer(C.fg_next_port_id))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_go_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_go_method_async))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_native_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_native_method_async))
	ptr ^= uintptr(unsafe.Pointer(C.fg_init_method_handle))
}

func unpackFgPacket(packet C.FgPacket) (method string, data []byte) {
	method = mapFgDataToString(packet.method)
	data = mapFgDataToBytes(packet.data)
	return
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
	len := C.int(len(from))
	return C.FgData{
		data: data,
		len:  len,
	}
}

func mapFgDataToBytes(from C.FgData) []byte {
	if from.data == nil {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(from.data), C.int(from.len))
}

func freeFgData(value C.FgData) {
	if value.data != nil {
		C.free(value.data)
		value.data = nil
	}
}
