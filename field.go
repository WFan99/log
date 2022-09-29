package log

import "fmt"

type FieldType uint8

const (
	StringType FieldType = iota
	IntType
)

type Field struct {
	fieldType FieldType
	key       string
	strValue  string
	numValue  int64
}

func String(key string, value string) Field {
	return Field{fieldType: StringType, key: key, strValue: value}
}

func Int(key string, value int) Field {
	return Field{fieldType: IntType, key: key, numValue: int64(value)}
}

func ErrObj(err error) Field {
	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}
	return Field{fieldType: StringType, key: "error", strValue: errStr}
}

func (field Field) AddTo(enc Encoder) {
	switch field.fieldType {
	case StringType:
		enc.AddString(field.key, field.strValue)
	case IntType:
		enc.AddInt(field.key, int(field.numValue))
	default:
		panic(fmt.Errorf("unkown field type %d", field.fieldType))
	}
}
