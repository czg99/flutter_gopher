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
	var method string
	if packet.method_len > 0 {
		method = C.GoStringN(packet.method, packet.method_len)
	} else {
		method = C.GoString(packet.method)
	}

	var data []byte
	if packet.data != nil {
		data = C.GoBytes(unsafe.Pointer(packet.data), C.int(packet.data_len))
		C.free(unsafe.Pointer(packet.data))
	}

	result := callGoMethod(method, data)

	var c_result unsafe.Pointer = nil
	if result != nil {
		c_result = C.CBytes(result)
	}
	c_result_len := C.int(len(result))
	return C.FgPacket{
		id:         packet.id,
		method:     packet.method,
		method_len: packet.method_len,
		data:       c_result,
		data_len:   c_result_len,
	}
}

//export fg_call_method_async
func fg_call_method_async(packet C.FgPacket) {
	go func() {
		result := fg_call_method(packet)
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
}
