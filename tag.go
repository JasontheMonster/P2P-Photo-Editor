package main

import (
	"fmt"
)

type Tag struct {
	Id			int		`json:"id"`
	Time_stamp	int 	`json:"time_stamp"`
}

func createTag(id int, ts int) Tag {
	return Tag{Id: id, Time_stamp: ts}
}

func (this *Tag) compareTo(other Tag) int {
	if this.Time_stamp != other.Time_stamp {
		return this.Time_stamp - other.Time_stamp
	} else {
		return this.Id - other.Id
	}
}

func (n *Node) updateTag(tag Tag) {
	var msg Message
    tmp := n.tag.compareTo(tag)
	if tmp < 0 {
		msg = n.createMessage(ACCEPT, "i need update", make(map[int]MemListEntry))
	} else if tmp > 0 {
		msg = n.createMessage(DECLINE, "i am newer", make(map[int]MemListEntry))
	} else {
        msg = n.createMessage(ACCEPT, "up to date", make(map[int]MemListEntry))
    }
	send(n.mem_list[tag.Id].Addr, msg)
}

type Message struct {
	// INVITE, PUBLIC, HEARTBEAT
    Kind        int 			       `json:"kind"`
    Ety         Entry 			       `json:"ety"`
    Tagval 		Tag 			       `json:"tagval"`
    Mem_list    map[int]MemListEntry   `json:"mem_list"`
    //QUIT          bool   
}

// function to create message
func (n *Node) createMessage(Kind int, info string, mem_list map[int]MemListEntry) Message {
    var msg Message
    msg.Kind = Kind
    msg.Ety = Entry{Time_stamp: n.tag.Time_stamp, Msg: info}
    msg.Tagval = createTag(n.ID, n.tag.Time_stamp)
    msg.Mem_list = mem_list
    return msg
}

//function to handle message
func (n *Node) handleMsg(msg Message){
    // fmt.Println(msg)
    switch msg.Kind {
    	case INVITE:
            fmt.Println(msg)
            n.tag.Time_stamp = msg.Tagval.Time_stamp
            n.log = initLog(msg.Tagval.Time_stamp)
            n.joinGroup(msg.Mem_list)
        case PUBLIC:
            n.checkPeers(msg.Mem_list)
            n.updateTag(msg.Tagval)
            fmt.Printf("Recved: %s\n", msg.Ety.Msg)
        case HEARTBEAT:
            n.checkPeers(msg.Mem_list)
            n.updateTag(msg.Tagval)
            fmt.Println("heartbeat", msg.Mem_list)
        case ACCEPT:
        	fmt.Printf("\tAccepted by %d\n", msg.Tagval.Id)
        case DECLINE:
            n.updateTag(msg.Tagval)
    }
}