package utils

import (
	"testing"
	"unsafe"
)

func Test_ToBytes(t *testing.T) {

	obj := SimpleMeta(10, 2, 3)
	t.Logf("%+v", obj)

	buf := (*[valEnd]byte)(unsafe.Pointer(obj))[:valEnd]
	buf[keyEnd] = byte(9)

	t.Logf("%+v", obj)
}

func Test_GetBytes(t *testing.T) {
	obj := &meta{
		HeaderSize: 10,
		KeySize:    2,
		ValSize:    3,
	}
	bytes := (*[valEnd]byte)(unsafe.Pointer(obj))[:valEnd]

	t.Logf("%+v", bytes)
}
