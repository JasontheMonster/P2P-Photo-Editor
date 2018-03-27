package main

// import (
// 	"math"
// 	"log"
// 	"time"
// )

type Tag struct {
	id			int			
	time_stamp	int
}

func createTag(id int, ts int) Tag {
	return Tag{id: id, time_stamp: ts}
}

func (this *Tag) smaller(other Tag) bool {
	if this.time_stamp < other.time_stamp {
		return true
	} else if this.time_stamp > other.time_stamp {
		return false
	} else {
		return this.id < other.id
	}
}

func (n *Node) updateTag(tag Tag) {
	var msg Message
	if n.tag.smaller(tag) {
		msg = n.createMessage(ACCEPT, "i need update", make(map[int]string))
	} else {
		msg = n.createMessage(DECLINE, "i am newer", make(map[int]string))
	}
	send(n.conn_list[tag.id], msg)
}

type Message struct {
    Kind        int  // INVITE, PUBLIC, HEARTBEAT
    ety         Entry
    tag 		Tag
    mem_list    map[int]string 
    //QUIT          bool   
}

func (n *Node) createMessage(Kind int, info string, mem_list map[int]string) Message {
    var msg Message
    msg.Kind = Kind
    msg.ety = Entry{time_stamp: n.tag.time_stamp, msg: info}
    msg.tag = createTag(n.ID, n.tag.time_stamp)
    msg.mem_list = mem_list
    return msg
}
