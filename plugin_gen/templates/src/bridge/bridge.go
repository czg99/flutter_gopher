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

//export fg_packet_loop
func fg_packet_loop() C.FgPacket {
	select {
	case result := <-packetChan:
		return result
	case <-time.After(time.Second):
		return fg_empty_packet()
	}
}

//export fg_call_method
func fg_call_method(packet C.FgPacket) C.FgPacket {
	method, data := fgPacketToGo(packet)
	if packet.data != nil {
		C.free(unsafe.Pointer(packet.data))
	}
	result := callGoMethod(method, data)
	return fgPacketFromGo(packet, result)
}

//export fg_call_method_async
func fg_call_method_async(packet C.FgPacket) {
	go func() {
		result := fg_call_method(packet)
		packetChan <- result
	}()
}

//export fg_call_native_method
func fg_call_native_method(packet C.FgPacket) C.FgPacket {
	method, data := fgPacketToGo(packet)
	if packet.data != nil {
		C.free(unsafe.Pointer(packet.data))
	}
	result := CallMethod(method, data)
	return fgPacketFromGo(packet, result)
}

//export fg_call_native_method_async
func fg_call_native_method_async(packet C.FgPacket) {
	go func() {
		result := fg_call_native_method(packet)
		packetChan <- result
	}()
}

//export enforce_binding
func enforce_binding() {
	var ptr uintptr
	ptr ^= uintptr(unsafe.Pointer(C.fg_next_port_id))
	ptr ^= uintptr(unsafe.Pointer(C.fg_empty_packet))
	ptr ^= uintptr(unsafe.Pointer(C.fg_packet_loop))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_method_async))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_native_method))
	ptr ^= uintptr(unsafe.Pointer(C.fg_call_native_method_async))
}

func fgPacketToGo(packet C.FgPacket) (method string, data []byte) {
	if packet.method != nil {
		if packet.method_len > 0 {
			method = C.GoStringN(packet.method, packet.method_len)
		} else {
			method = C.GoString(packet.method)
		}
	}
	if packet.data != nil {
		data = C.GoBytes(unsafe.Pointer(packet.data), C.int(packet.data_len))
	}
	return
}

func fgPacketFromGo(srcPacket C.FgPacket, data []byte) C.FgPacket {
	var cData unsafe.Pointer
	cDataLen := C.int(len(data))
	if data != nil {
		cData = C.CBytes(data)
	}
	return C.FgPacket{
		id:         srcPacket.id,
		method:     srcPacket.method,
		method_len: srcPacket.method_len,
		data:       cData,
		data_len:   cDataLen,
	}
}
