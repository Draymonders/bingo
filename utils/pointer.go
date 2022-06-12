package utils

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

/*
unsafe.Pointer 使用
*/
const (
	metaSize  = 16
	headerEnd = 8
	keyEnd    = 12
	valEnd    = 16
)

var (
	byteSerializer = binary.LittleEndian
)

type meta struct {
	HeaderSize int64
	KeySize    int32
	ValSize    int32
}

func SimpleMeta(headerSize int64, keySize int32, ValSize int) *meta {
	buf := make([]byte, 0, metaSize)

	buf = append(buf, ToBytes(headerSize)...)
	buf = append(buf, ToBytes(keySize)...)
	buf = append(buf, ToBytes(ValSize)...)

	return (*meta)(unsafe.Pointer(&buf[0]))
}

func ToBytes(val interface{}) (res []byte) {
	if val == nil {
		return res
	}
	switch val.(type) {
	case int64:
		v := val.(int64)
		res = make([]byte, 8)
		byteSerializer.PutUint64(res, uint64(v))
	case int:
		v := val.(int)
		res = make([]byte, 8)
		byteSerializer.PutUint32(res, uint32(v))
	case int32:
		v := val.(int32)
		res = make([]byte, 4)
		byteSerializer.PutUint32(res, uint32(v))
	default:
		panic(fmt.Sprintf("val: %v not support", val))
	}
	return
}

func GetBytes(bs []byte) interface{} {
	if len(bs) == 0 {
		return 0
	}
	if len(bs) == 4 {
		return byteSerializer.Uint32(bs)
	}
	if len(bs) == 8 {
		return byteSerializer.Uint64(bs)
	}
	return 0
}
