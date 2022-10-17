package log

import (
	"fmt"
	"math"
	"os"
	"time"
)

type FieldType uint8

const (
	// BinaryType indicates that the field carries an opaque binary blob.
	BinaryType FieldType = iota
	// BoolType indicates that the field carries a bool.
	BoolType
	// ByteStringType indicates that the field carries UTF-8 encoded bytes.
	ByteStringType
	// DurationType indicates that the field carries a time.Duration.
	DurationType
	// Float64Type indicates that the field carries a float64.
	Float64Type
	// Float32Type indicates that the field carries a float32.
	Float32Type
	// Int64Type indicates that the field carries an int64.
	Int64Type
	// Int32Type indicates that the field carries an int32.
	Int32Type
	// Int16Type indicates that the field carries an int16.
	Int16Type
	// Int8Type indicates that the field carries an int8.
	Int8Type
	// StringType indicates that the field carries a String.
	StringType
	// TimeType indicates that the field carries a time.Time that is
	// representable by a UnixNano() stored as an int64.
	TimeType
	// TimeFullType indicates that the field carries a time.Time stored as-is.
	TimeFullType
	// Uint64Type indicates that the field carries a uint64.
	Uint64Type
	// Uint32Type indicates that the field carries a uint32.
	Uint32Type
	// Uint16Type indicates that the field carries a uint16.
	Uint16Type
	// Uint8Type indicates that the field carries a uint8.
	Uint8Type
	// ReflectType indicates that the field carries an interface{}, which should
	// be serialized using reflection.
	ReflectType
)

var (
	_minTimeInt64 = time.Unix(0, math.MinInt64)
	_maxTimeInt64 = time.Unix(0, math.MaxInt64)
)

type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	String    string
	Interface interface{}
}

// nilField returns a field which will marshal explicitly as nil. See motivation
// in https://github.com/uber-go/zap/issues/753 . If we ever make breaking
// changes and add zapcore.NilType and zapcore.ObjectEncoder.AddNil, the
// implementation here should be changed to reflect that.
func nilField(key string) Field { return Reflect(key, nil) }

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) Field {
	return Field{Key: key, Type: BinaryType, Interface: val}
}

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) Field {
	var ival int64
	if val {
		ival = 1
	}
	return Field{Key: key, Type: BoolType, Integer: ival}
}

// Boolp constructs a field that carries a *bool. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Boolp(key string, val *bool) Field {
	if val == nil {
		return nilField(key)
	}
	return Bool(key, *val)
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) Field {
	return Field{Key: key, Type: ByteStringType, Interface: val}
}

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: Float64Type, Integer: int64(math.Float64bits(val))}
}

// Float64p constructs a field that carries a *float64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Float64p(key string, val *float64) Field {
	if val == nil {
		return nilField(key)
	}
	return Float64(key, *val)
}

// Float32 constructs a field that carries a float32. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float32(key string, val float32) Field {
	return Field{Key: key, Type: Float32Type, Integer: int64(math.Float32bits(val))}
}

// Float32p constructs a field that carries a *float32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Float32p(key string, val *float32) Field {
	if val == nil {
		return nilField(key)
	}
	return Float32(key, *val)
}

// Int constructs a field with the given key and value.
func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

// Intp constructs a field that carries a *int. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Intp(key string, val *int) Field {
	if val == nil {
		return nilField(key)
	}
	return Int(key, *val)
}

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: Int64Type, Integer: val}
}

// Int64p constructs a field that carries a *int64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int64p(key string, val *int64) Field {
	if val == nil {
		return nilField(key)
	}
	return Int64(key, *val)
}

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) Field {
	return Field{Key: key, Type: Int32Type, Integer: int64(val)}
}

// Int32p constructs a field that carries a *int32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int32p(key string, val *int32) Field {
	if val == nil {
		return nilField(key)
	}
	return Int32(key, *val)
}

// Int16 constructs a field with the given key and value.
func Int16(key string, val int16) Field {
	return Field{Key: key, Type: Int16Type, Integer: int64(val)}
}

// Int16p constructs a field that carries a *int16. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int16p(key string, val *int16) Field {
	if val == nil {
		return nilField(key)
	}
	return Int16(key, *val)
}

// Int8 constructs a field with the given key and value.
func Int8(key string, val int8) Field {
	return Field{Key: key, Type: Int8Type, Integer: int64(val)}
}

// Int8p constructs a field that carries a *int8. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Int8p(key string, val *int8) Field {
	if val == nil {
		return nilField(key)
	}
	return Int8(key, *val)
}

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return Field{Key: key, Type: StringType, String: val}
}

// Stringp constructs a field that carries a *String. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Stringp(key string, val *string) Field {
	if val == nil {
		return nilField(key)
	}
	return String(key, *val)
}

// Uint constructs a field with the given key and value.
func Uint(key string, val uint) Field {
	return Uint64(key, uint64(val))
}

// Uintp constructs a field that carries a *uint. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uintp(key string, val *uint) Field {
	if val == nil {
		return nilField(key)
	}
	return Uint(key, *val)
}

// Uint64 constructs a field with the given key and value.
func Uint64(key string, val uint64) Field {
	return Field{Key: key, Type: Uint64Type, Integer: int64(val)}
}

// Uint64p constructs a field that carries a *uint64. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint64p(key string, val *uint64) Field {
	if val == nil {
		return nilField(key)
	}
	return Uint64(key, *val)
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) Field {
	return Field{Key: key, Type: Uint32Type, Integer: int64(val)}
}

// Uint32p constructs a field that carries a *uint32. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint32p(key string, val *uint32) Field {
	if val == nil {
		return nilField(key)
	}
	return Uint32(key, *val)
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) Field {
	return Field{Key: key, Type: Uint16Type, Integer: int64(val)}
}

// Uint16p constructs a field that carries a *uint16. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint16p(key string, val *uint16) Field {
	if val == nil {
		return nilField(key)
	}
	return Uint16(key, *val)
}

// Uint8 constructs a field with the given key and value.
func Uint8(key string, val uint8) Field {
	return Field{Key: key, Type: Uint8Type, Integer: int64(val)}
}

// Uint8p constructs a field that carries a *uint8. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Uint8p(key string, val *uint8) Field {
	if val == nil {
		return nilField(key)
	}
	return Uint8(key, *val)
}

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]String to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) Field {
	return Field{Key: key, Type: ReflectType, Interface: val}
}

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) Field {
	if val.Before(_minTimeInt64) || val.After(_maxTimeInt64) {
		return Field{Key: key, Type: TimeFullType, Interface: val}
	}
	return Field{Key: key, Type: TimeType, Integer: val.UnixNano(), Interface: val.Location()}
}

// Timep constructs a field that carries a *time.Time. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Timep(key string, val *time.Time) Field {
	if val == nil {
		return nilField(key)
	}
	return Time(key, *val)
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: DurationType, Integer: int64(val)}
}

// Durationp constructs a field that carries a *time.Duration. The returned Field will safely
// and explicitly represent `nil` when appropriate.
func Durationp(key string, val *time.Duration) Field {
	if val == nil {
		return nilField(key)
	}
	return Duration(key, *val)
}

func ErrObj(err error) Field {
	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}
	return Field{Type: StringType, Key: "error", String: errStr}
}

func (f Field) AddTo(enc Encoder) {
	var err error

	switch f.Type {
	case BinaryType:
		enc.AddBinary(f.Key, f.Interface.([]byte))
	case BoolType:
		enc.AddBool(f.Key, f.Integer == 1)
	case ByteStringType:
		enc.AddByteString(f.Key, f.Interface.([]byte))
	case DurationType:
		enc.AddDuration(f.Key, time.Duration(f.Integer))
	case Float64Type:
		enc.AddFloat64(f.Key, math.Float64frombits(uint64(f.Integer)))
	case Float32Type:
		enc.AddFloat32(f.Key, math.Float32frombits(uint32(f.Integer)))
	case Int64Type:
		enc.AddInt64(f.Key, f.Integer)
	case Int32Type:
		enc.AddInt32(f.Key, int32(f.Integer))
	case Int16Type:
		enc.AddInt16(f.Key, int16(f.Integer))
	case Int8Type:
		enc.AddInt8(f.Key, int8(f.Integer))
	case StringType:
		enc.AddString(f.Key, f.String)
	case TimeType:
		if f.Interface != nil {
			enc.AddTime(f.Key, time.Unix(0, f.Integer).In(f.Interface.(*time.Location)))
		} else {
			// Fall back to UTC if location is nil.
			enc.AddTime(f.Key, time.Unix(0, f.Integer))
		}
	case TimeFullType:
		enc.AddTime(f.Key, f.Interface.(time.Time))
	case Uint64Type:
		enc.AddUint64(f.Key, uint64(f.Integer))
	case Uint32Type:
		enc.AddUint32(f.Key, uint32(f.Integer))
	case Uint16Type:
		enc.AddUint16(f.Key, uint16(f.Integer))
	case Uint8Type:
		enc.AddUint8(f.Key, uint8(f.Integer))
	case ReflectType:
		err = enc.AddReflected(f.Key, f.Interface)
	default:
		panic(fmt.Sprintf("unknown field type: %v", f))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError%v", f.Key, err.Error())
	}
}
