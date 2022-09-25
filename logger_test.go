package coolog

import (
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	logger := NewDefaultLogger()
	logger.Debug("hello world", Int("num", 50), String("name", "wangfan"))
}

func TestAsyncLog(t *testing.T) {
	logger := NewDefaultLogger()
	logger.SetStdout(false)
	allLogConf := AsyncLoggerConfig{path: "coolog.all.log", rollTime: 3600, rollSize: 512 * 1024 * 1024, logLevel: DebugLevel}
	warnLogConf := AsyncLoggerConfig{path: "coolog.warn.log", rollTime: 3600, rollSize: 512 * 1024 * 1024, logLevel: WarnLevel}
	logger.StartAsyncLog(allLogConf, warnLogConf)

	logger.Debug("debug log", Int("num", 50), String("name", "wangfan"), Int("age", 23))
	logger.Warn("hello world", Int("num", 50), String("name", "wangfan"), Int("age", 23))
	logger.StopAsyncLog()
}
