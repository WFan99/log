package log

import "time"

type EntryCaller struct {
	Defined bool
	File    string
	Line    int
}

type Entry struct {
	Level   Level
	Time    time.Time
	Message string
	Caller  EntryCaller
}
