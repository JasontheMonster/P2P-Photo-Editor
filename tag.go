package main

import (
	"encoding/json"
// 	"math"
// 	"log"
// 	"time"
)

type Tag struct {
	Id			int			
	Time_stamp	int
}

func createTag(id int, ts int) Tag {
	return Tag{Id: id, Time_stamp: ts}
}

func (this *Tag) smaller(other Tag) bool {
	if this.Time_stamp < other.Time_stamp {
		return true
	} else if this.Time_stamp > other.Time_stamp {
		return false
	} else {
		return this.Id < other.Id
	}
}

func (n *Node) updateTag(tag Tag) {
	var msg Message
	if n.tag.smaller(tag) {
		msg = n.createMessage(ACCEPT, "i need update", make(map[int]string))
	} else {
		msg = n.createMessage(DECLINE, "i am newer", make(map[int]string))
	}
		json.NewEncoder(n.conn_list[tag.Id]).Encode(msg)
}

type Message struct {
    Kind        int  // INVITE, PUBLIC, HEARTBEAT
    Ety         Entry
    Tagval 		Tag
    Mem_list    map[int]string 
    //QUIT          bool   
}

func (n *Node) createMessage(Kind int, info string, mem_list map[int]string) Message {
    var msg Message
    msg.Kind = Kind
    msg.Ety = Entry{Time_stamp: n.tag.Time_stamp, Msg: info}
    msg.Tagval = createTag(n.ID, n.tag.Time_stamp)
    msg.Mem_list = mem_list
    return msg
}
