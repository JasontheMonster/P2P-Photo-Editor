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

// initialize log 
func initLog(ts int) Log {
	var ety []Entry
	return Log{Last_applied: 0, Time_stamp: ts, Entries: ety}
}

// append entry to log and update time
func (l *Log) append(ety Entry){
	mutex.Lock()
	l.Time_stamp += 1
	l.Entries = append(l.Entries, ety)
	mutex.Unlock()
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

// send commited log entries to front end
func (n *Node) applyLog() {
	fmt.Printf("nts: %d, la: %d, ts: %d\n", n.tag.Time_stamp, n.log.Last_applied, n.log.Time_stamp)
	//lock the critical session of incrementing log's lastapplied var
	mutex.Lock()
	// sequentially send log entry to front end
	for i := n.log.Last_applied; i < len(n.log.Entries); i++ {
		// take snapshot every 15 log entries
		if i % MAX_LOG_ENTRY == 0 {
			n.sendToFront("snapshot")
		}
		n.sendToFront(n.log.Entries[i].Msg)
		fmt.Println(n.log.Entries[i])
		n.log.Last_applied += 1
	}
	mutex.Unlock()
}