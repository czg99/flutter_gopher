package bridge

import "log"

// MethodHandle 处理函数
type MethodHandle func(method string, data []byte) []byte

var goMethodHandle MethodHandle = nil

// InitMethodHandle 初始化处理函数
func InitMethodHandle(handle MethodHandle) {
	goMethodHandle = handle
}

// callGoMethod 调用 go 的处理函数
func callGoMethod(method string, data []byte) []byte {
	if goMethodHandle == nil {
		log.Println("go method handle not init")
		return nil
	}
	return goMethodHandle(method, data)
}
