package main

// type log entry
type Entry struct {
	Time_stamp	int 	`json:"time_stamp"`
	Msg			string 	`json:"msg"`
}

type Log struct {
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
	return Log{Time_stamp: ts, Entries: ety}
}

// append entry to log and update time
func (l *Log) append(ety Entry){
	l.Time_stamp = ety.Time_stamp
	l.Entries = append(l.Entries, ety)
}

// update logs with recved updatelog
func (l *Log) updateLog(etys []Entry) {
	for _, ety := range etys {
		if l.Time_stamp + 1 == ety.Time_stamp {
			l.append(ety)
		}
	}
}

// if older than incoming heartbeat, send update request
func (n *Node) checkLog(tag Tag) {
	if n.tag.compareTo(tag) < 0 {
		req := n.createMessage(UPDATEREQUEST, "", make(map[int]MemListEntry))
        send(n.mem_list[tag.ID].Addr, req)
	}
}