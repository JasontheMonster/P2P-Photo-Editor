package main

type Entry struct {
	Time_stamp	int 	`json:"time_stamp"`
	Msg			string 	`json:"msg"`
}

type Log struct {
	time_stamp	int
	entries		[]Entry
}

type HoldBackEty struct {
	Ety 	Entry
	Time 	int64
}
//initialize log 
func initLog(ts int) Log {
	var ety []Entry
	return Log{time_stamp: ts, entries: ety}
}

func (l *Log)append(ety Entry){
	l.time_stamp = ety.Time_stamp
	l.entries = append(l.entries, ety)
}

func (l *Log) updateLog(etys []Entry) {
	for _, ety := range etys {
		if l.time_stamp + 1 == ety.Time_stamp {
			l.append(ety)
		}
	}
}

func (n *Node) checkLog(tag Tag) {
	if n.tag.compareTo(tag) < 0 {
		req := n.createMessage(UPDATEREQUEST, "", make(map[int]MemListEntry))
        send(n.mem_list[tag.Id].Addr, req)
	}
}