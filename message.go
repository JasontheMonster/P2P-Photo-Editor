package main

import (
	"fmt"
)

type Message struct {
    Kind        int 			       `json:"kind"`
    Ety         Entry 			       `json:"ety"`
    Tag 		Tag 			       `json:"tagval"`
    Mem_list    map[int]MemListEntry   `json:"mem_list"`
    UpdateInfo  []Entry                `json:"updateinfo"`
  
}

// function to create message
func (n *Node) createMessage(Kind int, info string, mem_list map[int]MemListEntry) Message {
    var msg Message
    msg.Kind = Kind
    msg.Ety = Entry{Time_stamp: n.tag.Time_stamp, Msg: info}
    msg.Tag = createTag(n.ID, n.tag.Time_stamp)
    msg.Mem_list = mem_list
    msg.UpdateInfo = make([]Entry, 1)
    return msg
}

// create log file message
func (n *Node) createMessageWithLog(Kind int, info string, updateinfo[]Entry) Message{
    msg := n.createMessage(Kind, info, make(map[int]MemListEntry))
    msg.UpdateInfo = updateinfo
    return msg
}

func (n *Node) createUpdateRequest() Message {
    msg := n.createMessage(UPDATEREQUEST, "", make(map[int]MemListEntry))
    msg.Tag.Time_stamp = n.log.Time_stamp
    return msg
}

// data message's tag is current tag + 1
func (n *Node) createDataMessage(Kind int, info string) Message {
    msg := n.createMessage(Kind, info, make(map[int]MemListEntry))
    msg.Ety.Time_stamp += 1
    return msg
}

func (n *Node) commit(msg Message) {
    n.voted = false
    if (n.tag.compareTo(msg.Tag) == -1) {
        n.log.append(n.holdBack.Ety)
    }
    n.tag.Time_stamp = msg.Tag.Time_stamp 
    fmt.Printf("Commited: %s, %d\n", n.holdBack.Ety.Msg, n.holdBack.Ety.Time_stamp)
    n.applyLog()
}

//function to handle message
func (n *Node) handleMsg(msg Message){
    switch msg.Kind {
    	case INVITE:
            fmt.Println("\tAccepted invitation.")
            n.tag.Time_stamp = msg.Tag.Time_stamp
            n.log = initLog(0)
            targetId := msg.Tag.ID
            n.joinGroup(msg.Mem_list, targetId)
            if n.HasImage{
                req := n.createUpdateRequest()
                send(n.mem_list[msg.Tag.ID].Addr, req)
            }
        case PUBLIC:
            n.checkPeers(msg.Mem_list)
            n.updateTag(msg)
            fmt.Printf("\tRecved: %s\n", msg.Ety.Msg)
        case HEARTBEAT:
            if n.HasImage {
                n.checkPeers(msg.Mem_list)
                n.checkLog(msg.Tag)
            }
            fmt.Println("\theartbeat")
        case ACCEPT:
        	fmt.Printf("\tInvite accepted by %d\n", msg.Tag.ID)
            n.mem_list[msg.Tag.ID] = msg.Mem_list[msg.Tag.ID]
        // case DECLINE:
        //     n.updateTag(msg)
        case ACK:
            if ack, isIn := chans[msg.Ety.Time_stamp]; isIn{
                ack <- (msg.Ety.Msg == "agreed" )
            }
        case COMMIT:
            n.commit(msg)
        case UPDATEINFO:
            fmt.Println("recv update")
            n.log.updateLog(msg.UpdateInfo)
            n.applyLog()
            // fmt.Printf("receive update info from %d\n", msg.Tagval.Id)
        case UPDATEREQUEST:
            fmt.Println("recv request")
            n.sendUpdate(msg.Tag)
            // fmt.Printf("sent update info to %d\n", msg.Tagval.Id)
    }
}