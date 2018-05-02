//Tag has information about node's logical tiimestamp and how to compare them
package main

import (
    "time"
    // "fmt"
)

type Tag struct {
    // id of the sender
	ID			int		`json:"id"`
    // logical timestamp
	Time_stamp	int 	`json:"time_stamp"`
}

// construct tag by id and timestamp
func createTag(id int, ts int) Tag {
	return Tag{ID: id, Time_stamp: ts}
}

// compare tag by timestamp
func (this *Tag) compareTo(other Tag) int {
	return this.Time_stamp - other.Time_stamp
}

// repond to an incoming public message
func (n *Node) updateTag(msg Message) {
    mutex.Lock()
	var rep Message
    tmp := n.tag.compareTo(msg.Tag)
    if n.voted || tmp > 0 { // if already voted or msg has older tag, decline the message
        rep = n.createMessage(ACK, "fuck ya", make(map[int]MemListEntry))
    } else { // if not voted and msg has newer tag, accept the message
        n.voted = true
        n.holdBack = time.Now().UnixNano()
		rep = n.createDataMessage(ACK, "agreed")
        if tmp < -1 { // if self is not up to date, request for update
            req := n.createUpdateRequest()
            send(n.mem_list[msg.Tag.ID].Addr, req)
        }
	}
    // send the response
	send(n.mem_list[msg.Tag.ID].Addr, rep)
    mutex.Unlock()
}