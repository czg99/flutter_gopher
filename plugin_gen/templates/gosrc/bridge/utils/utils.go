package utils

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

func CValueToPtr[T any](value T) *T {
	size := unsafe.Sizeof(value)
	data := C.malloc(C.size_t(size))
	*(*T)(data) = value
	return (*T)(data)
}
