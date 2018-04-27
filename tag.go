package main

import (
    "time"
)

type Tag struct {
	ID			int		`json:"id"`
	Time_stamp	int 	`json:"time_stamp"`
}

func createTag(id int, ts int) Tag {
	return Tag{ID: id, Time_stamp: ts}
}

func (this *Tag) compareTo(other Tag) int {
	return this.Time_stamp - other.Time_stamp
}

func (n *Node) updateTag(msg Message) {
	var rep Message
    tmp := n.tag.compareTo(msg.Tag)
    if n.voted || tmp > 0 {
        rep = n.createMessage(ACK, "fuck ya", make(map[int]MemListEntry))
    } else {
        n.voted = true
        n.holdBack = HoldBackEty{Ety: msg.Ety, Time: time.Now().UnixNano()}
		rep = n.createDataMessage(ACK, "agreed")
        if tmp < 0 {
            req := n.createUpdateRequest()
            send(n.mem_list[msg.Tag.ID].Addr, req)
        }
	}
	send(n.mem_list[msg.Tag.ID].Addr, rep)
}