package photoEditor

import (
	"math"
	"log"
	"time"
)

type Entry struct {
	time_stamp	int
	log_entry	string
}

type Log struct {
	time_stamp	int
	entries		[]Entry
}

func (l *Log)append(log_entry string){
	
}
