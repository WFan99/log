package log

import (
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	Debug("hello world", Int("num", 50), String("name", "wangfan"))
}

func TestAsyncLog(t *testing.T) {
	logger := NewDefaultLogger()
	logger.SetStdout(false)
	allLogConf := AsyncLoggerConfig{path: "test.all.log", rollTime: 3600, rollSize: 512 * 1024 * 1024, logLevel: DebugLevel}
	warnLogConf := AsyncLoggerConfig{path: "test.warn.log", rollTime: 3600, rollSize: 512 * 1024 * 1024, logLevel: WarnLevel}
	logger.StartAsyncLog(allLogConf, warnLogConf)

	logger.Debug("debug log", Int("num", 50), String("name", "wangfan"), Int("age", 23))
	logger.Warn("hello world", Int("num", 50), String("name", "wangfan"), Int("age", 23))
	logger.StopAsyncLog()
}
