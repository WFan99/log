package log

import (
	"os"
)

type Logger struct {
	stdout       bool
	logLevel     AtomicLevel
	asyncLoggers map[string]*AsyncLogger
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
	enc := &JsonEncoder{}
	bytes := enc.Encode(msg, level, fields...)
	if logger.stdout {
		os.Stdout.Write(bytes)
	}
	for _, al := range logger.asyncLoggers {
		al.appendContent(&logContent{level: level, bytes: bytes})
	}

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
