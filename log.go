package main

import "fmt"

// type log entry
type Entry struct {
	Time_stamp	int 	`json:"time_stamp"`
	Msg			string 	`json:"msg"`
}

type Log struct {
	Last_applied 	int
	Time_stamp		int
	Entries			[]Entry
}

type HoldBackEty struct {
	Ety 	Entry
	Time 	int64
}

//initialize log 
func initLog(ts int) Log {
	var ety []Entry
	return Log{Last_applied: 0, Time_stamp: ts, Entries: ety}
}

// append entry to log and update time
func (l *Log) append(ety Entry){
	l.Time_stamp = ety.Time_stamp
	l.Entries = append(l.Entries, ety)
}

// update logs with recved updatelog
func (l *Log) updateLog(etys []Entry) {
	for _, ety := range etys {
		if l.Time_stamp == ety.Time_stamp - 1 {
			l.append(ety)
		}
	}
}

// if older than incoming heartbeat, send update request
func (n *Node) checkLog(tag Tag) {
	if n.log.Time_stamp < tag.Time_stamp {
		req := n.createUpdateRequest()
        send(n.mem_list[tag.ID].Addr, req)
	}
}

func (n *Node) applyLog() {
	fmt.Printf("nts: %d, la: %d, ts: %d\n", n.tag.Time_stamp, n.log.Last_applied, n.log.Time_stamp)
	for i := n.log.Last_applied; i < len(n.log.Entries); i++ {
		n.sendToFront(n.log.Entries[i].Msg)
		n.log.Last_applied += 1
	}
}