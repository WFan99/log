package log

import "time"

type ObjectEncoder interface {
	// Built-in types.
	AddBinary(key string, value []byte)     // for arbitrary bytes
	AddByteString(key string, value []byte) // for UTF-8 encoded bytes
	AddBool(key string, value bool)
	AddDuration(key string, value time.Duration)
	AddFloat64(key string, value float64)
	AddFloat32(key string, value float32)
	AddInt(key string, value int)
	AddInt64(key string, value int64)
	AddInt32(key string, value int32)
	AddInt16(key string, value int16)
	AddInt8(key string, value int8)
	AddString(key, value string)
	AddTime(key string, value time.Time)
	AddUint(key string, value uint)
	AddUint64(key string, value uint64)
	AddUint32(key string, value uint32)
	AddUint16(key string, value uint16)
	AddUint8(key string, value uint8)
	AddUintptr(key string, value uintptr)

	// AddReflected uses reflection to serialize arbitrary objects, so it can be
	// slow and allocation-heavy.
	AddReflected(key string, value interface{}) error
}

type Encoder interface {
	ObjectEncoder
	EncodeEntry(Entry, []Field) (*Buffer, error)
}
