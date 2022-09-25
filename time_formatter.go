package coolog

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

var gTimeFormatter TimeFormatter

type TimeFormatter struct {
	mtx               sync.Mutex
	lastEncodeTimeSec int64
	lastEncodeTimePre string
}

func (tf *TimeFormatter) getTimeString(t time.Time) string {
	sb := strings.Builder{}
	tf.mtx.Lock()
	if tf.lastEncodeTimeSec != t.Unix() {
		tf.lastEncodeTimePre = t.Format("2006-01-02 15:04:05")
	}
	sb.WriteString(tf.lastEncodeTimePre)
	tf.mtx.Unlock()
	sb.WriteByte('.')
	micro := t.Nanosecond() / 1000
	if micro < 100000 {
		sb.WriteByte('0')
	}
	if micro < 10000 {
		sb.WriteByte('0')
	}
	if micro < 1000 {
		sb.WriteByte('0')
	}
	if micro < 100 {
		sb.WriteByte('0')
	}
	if micro < 10 {
		sb.WriteByte('0')
	}
	sb.WriteString(strconv.Itoa(micro))
	return sb.String()
}
