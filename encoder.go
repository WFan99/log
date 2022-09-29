package log

type Encoder interface {
	Encode(msg string, level Level, fields ...Field) []byte
	AddString(key string, value string)
	AddInt(key string, value int)
}
