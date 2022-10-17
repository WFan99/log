package log

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

// For JSON-escaping; see jsonEncoder.safeAddString below.
const _hex = "0123456789abcdef"

type jsonEncoder struct {
	buf *Buffer
}

var JsonPool = sync.Pool{New: func() interface{} {
	return &jsonEncoder{}
}}

func getJsonEncoder() *jsonEncoder {
	je := JsonPool.Get().(*jsonEncoder)
	je.buf = smallBufferPool.Get()
	return je
}

func putJsonEncoder(je *jsonEncoder) {
	if je.buf != nil {
		je.buf.Free()
		je.buf = nil
	}
	JsonPool.Put(je)
}

func (enc *jsonEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *jsonEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *jsonEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *jsonEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *jsonEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *jsonEncoder) AddFloat32(key string, val float32) {
	enc.addKey(key)
	enc.AppendFloat32(val)
}

func (enc *jsonEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

var nullLiteralBytes = []byte("null")

func (enc *jsonEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *jsonEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *jsonEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

// AddReflected uses reflection to serialize arbitrary objects, so it can be
// slow and allocation-heavy.
func (enc *jsonEncoder) AddReflected(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	enc.addKey(key)
	enc.buf.Write(bytes)
	return nil
}

func (enc *jsonEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *jsonEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

// appendComplex appends the encoded form of the provided complex128 value.
// precision specifies the encoding precision for the real and imaginary
// components of the complex number.
func (enc *jsonEncoder) appendComplex(val complex128, precision int) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted String, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, precision)
	// If imaginary part is less than 0, minus (-) sign is added by default
	// by AppendFloat.
	if i >= 0 {
		enc.buf.AppendByte('+')
	}
	enc.buf.AppendFloat(i, precision)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) AppendDuration(val time.Duration) {
	enc.AppendInt64(int64(val))
}

func (enc *jsonEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *jsonEncoder) AppendString(val string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(val)
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) AppendTimeLayout(time time.Time, layout string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.buf.AppendTime(time, layout)
	enc.buf.AppendByte('"')
}

func (enc *jsonEncoder) AppendTime(val time.Time) {
	// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
	// output JSON valid.
	enc.AppendString(gTimeFormatter.getTimeString(val))
}

func (enc *jsonEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *jsonEncoder) AddInt(k string, v int)         { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt32(k string, v int32)     { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt16(k string, v int16)     { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt8(k string, v int8)       { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddUint(k string, v uint)       { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint32(k string, v uint32)   { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint16(k string, v uint16)   { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint8(k string, v uint8)     { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUintptr(k string, v uintptr) { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AppendFloat64(v float64)        { enc.appendFloat(v, 64) }
func (enc *jsonEncoder) AppendFloat32(v float32)        { enc.appendFloat(float64(v), 32) }
func (enc *jsonEncoder) AppendInt(v int)                { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt32(v int32)            { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt16(v int16)            { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt8(v int8)              { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendUint(v uint)              { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint32(v uint32)          { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint16(v uint16)          { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint8(v uint8)            { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUintptr(v uintptr)        { enc.AppendUint64(uint64(v)) }

func (enc *jsonEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddString(key)
	enc.buf.AppendByte('"')
	enc.buf.AppendByte(':')
}

func (enc *jsonEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
	}
}

func (enc *jsonEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a String and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *jsonEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(String(s)) for s []byte.
func (enc *jsonEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *jsonEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *jsonEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func (enc *jsonEncoder) EncodeEntry(entry Entry, fields []Field) (*Buffer, error) {
	enc.buf.AppendByte('{')
	enc.AddString("level", entry.Level.String())
	enc.AddTime("time", entry.Time)
	if entry.Caller.Defined {
		enc.AddString("file", entry.Caller.File+":"+strconv.Itoa(entry.Caller.Line))
	}
	enc.AddString("message", entry.Message)
	for _, field := range fields {
		field.AddTo(enc)
	}
	enc.buf.AppendByte('}')
	return enc.buf, nil
}
