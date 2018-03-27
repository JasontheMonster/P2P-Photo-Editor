package main

// import (
// 	"math"
// 	"log"
// 	"time"
// )

type Entry struct {
	Time_stamp	int
	Msg			string
}

type Log struct {
	time_stamp	int
	entries		[]Entry
}

func initLog(ts int) Log {
	var ety []Entry
	return Log{time_stamp: ts, entries: ety}
}

func (l *Log)append(ety Entry){
	l.time_stamp = ety.Time_stamp
	l.entries = append(l.entries, ety)
}
