package coolog

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type logContent struct {
	level Level
	bytes []byte
}

type AsyncLoggerConfig struct {
	path     string
	rollSize uint64 // byte
	rollTime int64  // second
	logLevel Level
}

type AsyncLogger struct {
	path           string
	rollSize       uint64 // byte
	rollTime       int64
	logLevel       AtomicLevel
	curWriteBytes  uint64
	bufferChan     chan []byte
	file           *os.File
	createFileTime time.Time
	wg             sync.WaitGroup
}

func newAsyncLogger(config AsyncLoggerConfig) *AsyncLogger {
	logger := &AsyncLogger{
		path:     config.path,
		rollSize: config.rollSize,
		rollTime: config.rollTime,
	}
	logger.logLevel.set(config.logLevel)
	return logger
}

func (al *AsyncLogger) start() {
	al.bufferChan = make(chan []byte)
	al.rollFile(true)
	al.wg = sync.WaitGroup{}
	al.wg.Add(1)
	go func() {
		running := true
		ticker := time.NewTicker(time.Second)
		for running {
			select {
			case buffer, ok := <-al.bufferChan:
				if !ok {
					running = false
					break
				}
				if al.file != nil {
					n, _ := al.file.Write(buffer)
					// 按大小归档
					al.curWriteBytes += uint64(n)
					if al.rollSize > 0 && al.curWriteBytes >= al.rollSize {
						al.rollFile(true)
					}
				}
			case <-ticker.C:
				// 按时间归档
				if al.rollTime > 0 {
					now := time.Now().Unix()
					if now%al.rollTime == 0 {
						al.rollFile(true)
					}
				}
			}
		}
		if al.file != nil {
			al.file.Sync()
		}
		al.rollFile(false)
		al.wg.Done()
	}()
}

func (al *AsyncLogger) rollFile(createFile bool) {
	if al.file != nil {
		al.file.Close()
		baseName := al.path + "." + al.createFileTime.Format("20060102.1504")
		finalName := baseName
		for i := 0; ; i++ {
			if i > 0 {
				finalName = baseName + "." + strconv.Itoa(i)
			}
			_, err := os.Stat(finalName)
			if os.IsNotExist(err) {
				break
			}
		}
		err := os.Rename(al.path, finalName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if createFile {
		file, err := os.OpenFile(al.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("open file failed %v\n", err)
			return
		}
		al.createFileTime = time.Now()
		al.file = file
	}
	al.curWriteBytes = 0
}

func (al *AsyncLogger) appendContent(content *logContent) {
	if content == nil {
		return
	}
	if al.logLevel.get().enableOutput(content.level) {
		al.bufferChan <- content.bytes
	}
}

func (al *AsyncLogger) stop() {
	close(al.bufferChan)
	al.wg.Wait()

}
