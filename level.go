package log

import "sync/atomic"

type Level uint8

const (
	DebugLevel Level = 1
	InfoLevel  Level = 2
	WarnLevel  Level = 3
	ErrorLevel Level = 4
)

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		panic("unknown log level")
	}
}

func (level Level) enableOutput(outputLevel Level) bool {
	return outputLevel >= level
}

type AtomicLevel struct {
	level uint32
}

func (al *AtomicLevel) get() Level {
	return Level(atomic.LoadUint32(&al.level))
}

func (al *AtomicLevel) set(level Level) {
	atomic.StoreUint32(&al.level, uint32(level))
}
