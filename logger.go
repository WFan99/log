package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

type Logger struct {
	stdout        bool
	logLevel      AtomicLevel
	addCallerSkip int
	asyncLoggers  map[string]*AsyncLogger
}

func NewDefaultLogger() *Logger {
	logger := &Logger{}
	logger.stdout = true
	logger.SetLogLevel(DebugLevel)
	logger.asyncLoggers = map[string]*AsyncLogger{}
	return logger
}

func (logger *Logger) GetLogLevel() Level {
	return logger.logLevel.get()
}

func (logger *Logger) SetLogLevel(level Level) {
	logger.logLevel.set(level)
}

func (logger *Logger) SetStdout(enable bool) {
	logger.stdout = enable
}

func (logger *Logger) StartAsyncLog(configs ...AsyncLoggerConfig) {
	for _, config := range configs {
		logger.asyncLoggers[config.path] = newAsyncLogger(config)
	}
	for _, al := range logger.asyncLoggers {
		al.start()
	}
}

func (logger *Logger) StopAsyncLog() {
	for _, asyncLogger := range logger.asyncLoggers {
		asyncLogger.stop()
	}
}

func (logger *Logger) output(msg string, level Level, fields ...Field) {
	if !logger.GetLogLevel().enableOutput(level) {
		return
	}
	entryCaller := EntryCaller{Defined: false}
	if _, file, line, ok := runtime.Caller(2 + logger.addCallerSkip); ok {
		entryCaller.Defined = true
		entryCaller.File = file
		entryCaller.Line = line
	}
	enc := getJsonEncoder()
	entry := Entry{Level: level, Time: time.Now(), Message: msg, Caller: entryCaller}
	buffer, err := enc.EncodeEntry(entry, fields)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v write error %v\n", entry.Time, err)
		return
	}
	if logger.stdout {
		os.Stdout.Write(buffer.Bytes())
	}
	for _, al := range logger.asyncLoggers {
		if al.logLevel.get().enableOutput(entry.Level) {
			al.appendBytes(buffer.Bytes())
		}
	}
	putJsonEncoder(enc)
	buffer.Free()
}

func (logger *Logger) Debug(msg string, fields ...Field) {
	logger.output(msg, DebugLevel, fields...)
}

func (logger *Logger) Info(msg string, fields ...Field) {
	logger.output(msg, InfoLevel, fields...)
}

func (logger *Logger) Warn(msg string, fields ...Field) {
	logger.output(msg, WarnLevel, fields...)
}

func (logger *Logger) Error(msg string, fields ...Field) {
	logger.output(msg, ErrorLevel, fields...)
}
