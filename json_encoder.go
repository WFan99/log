package log

import (
	"bytes"
	"strconv"
	"time"
)

type JsonEncoder struct {
	buffer bytes.Buffer
}

func (enc *JsonEncoder) Encode(msg string, level Level, fields ...Field) []byte {
	now := time.Now()
	enc.buffer.WriteByte('{')
	enc.appendKey("time")
	enc.appendString(gTimeFormatter.getTimeString(now))
	enc.buffer.WriteByte(',')
	enc.appendKey("level")
	enc.appendString(level.string())
	enc.buffer.WriteByte(',')
	enc.appendKey("message")
	enc.appendString(msg)

	for _, field := range fields {
		enc.buffer.WriteByte(',')
		field.AddTo(enc)
	}
	enc.buffer.WriteByte('}')
	enc.buffer.WriteByte('\n')
	return enc.buffer.Bytes()
}

func (enc *JsonEncoder) appendKey(key string) {
	enc.buffer.WriteByte('"')
	enc.buffer.WriteString(key)
	enc.buffer.WriteString(`":`)
}

func (enc *JsonEncoder) appendString(str string) {
	enc.buffer.WriteByte('"')
	enc.buffer.WriteString(str)
	enc.buffer.WriteByte('"')
}

func (enc *JsonEncoder) AddString(key string, value string) {
	enc.appendKey(key)
	enc.appendString(value)
}

func (enc *JsonEncoder) AddInt(key string, value int) {
	enc.appendKey(key)
	enc.appendString(strconv.Itoa(value))
}
