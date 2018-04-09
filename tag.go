package main

import (
	"fmt"
    "time"
)

type Tag struct {
	Id			int		`json:"id"`
	Time_stamp	int 	`json:"time_stamp"`
}

func createTag(id int, ts int) Tag {
	return Tag{Id: id, Time_stamp: ts}
}

func (this *Tag) compareTo(other Tag) int {
	return this.Time_stamp - other.Time_stamp
}

func (n *Node) updateTag(msg Message) {
	var rep Message
    tmp := n.tag.compareTo(msg.Tagval)
    if n.voted || tmp > 0 {
        rep = n.createMessage(ACK, "fuck ya", make(map[int]MemListEntry))
    } else {
        n.voted = true
        n.holdBack = HoldBackEty{Ety: msg.Ety, Time: time.Now().UnixNano()}
		rep = n.createDataMessage(ACK, "agreed", make(map[int]MemListEntry))
        if tmp < 0 {
            req := n.createMessage(UPDATEREQUEST, "", make(map[int]MemListEntry))
            send(n.mem_list[msg.Tagval.Id].Addr, req)
        }
	}
	send(n.mem_list[msg.Tagval.Id].Addr, rep)
}

type Message struct {
    Kind        int 			       `json:"kind"`
    Ety         Entry 			       `json:"ety"`
    Tagval 		Tag 			       `json:"tagval"`
    Mem_list    map[int]MemListEntry   `json:"mem_list"`
    UpdateInfo  []Entry                `json:"updateinfo"`

    //image file string
    Image       string                 `json:"image"`
    //QUIT          bool   
}


//create log file message

func (n *Node) createMessageWithLog(Kind int, info string, updateinfo[]Entry) Message{
    var msg Message
    msg.Kind = Kind
    msg.Ety = Entry{Time_stamp: n.tag.Time_stamp, Msg: info}
    msg.Tagval = createTag(n.ID, n.tag.Time_stamp)
    msg.Mem_list = make(map[int]MemListEntry)
    msg.UpdateInfo = updateinfo
    msg.Image = ""
    return msg
}
// function to create message
func (n *Node) createMessage(Kind int, info string, mem_list map[int]MemListEntry) Message {
    var msg Message
    msg.Kind = Kind
    msg.Ety = Entry{Time_stamp: n.tag.Time_stamp, Msg: info}
    msg.Tagval = createTag(n.ID, n.tag.Time_stamp)
    msg.Mem_list = mem_list
    msg.UpdateInfo = make([]Entry, 1)
    msg.Image = ""
    return msg
}

//function to send image
func (n *Node) createMessageWithImage(Kind int, info string, image string) Message {
    var msg Message
    msg.Kind = Kind
    msg.Ety = Entry{Time_stamp: n.tag.Time_stamp, Msg: info}
    msg.Tagval = createTag(n.ID, n.tag.Time_stamp)
    msg.Mem_list =  make(map[int]MemListEntry)
    msg.UpdateInfo = make([]Entry, 1)
    msg.Image = image
    return msg
}

// data message's tag is current tag + 1
func (n *Node) createDataMessage(Kind int, info string, mem_list map[int]MemListEntry) Message {
    msg := n.createMessage(Kind, info, mem_list)
    msg.Ety.Time_stamp += 1
    return msg
}

//function to handle message
func (n *Node) handleMsg(msg Message){
    switch msg.Kind {
    	case INVITE:
            fmt.Println("\tAccepted invitation.")
            n.tag.Time_stamp = msg.Tagval.Time_stamp
            n.log = initLog(msg.Tagval.Time_stamp)
            n.joinGroup(msg.Mem_list)
            // fmt.Println("image: ", msg.Image)
        //received 
        case PUBLIC:
            n.checkPeers(msg.Mem_list)
            n.updateTag(msg)
            fmt.Printf("\tRecved: %s\n", msg.Ety.Msg)
        case HEARTBEAT:
            n.checkPeers(msg.Mem_list)
            n.checkLog(msg.Tagval)
            // n.updateTag(msg)
            fmt.Println("\theartbeat")
        case ACCEPT:
        	fmt.Printf("\tInvite accepted by %d\n", msg.Tagval.Id)
            n.mem_list[msg.Tagval.Id] = msg.Mem_list[msg.Tagval.Id]
        // case DECLINE:
        //     n.updateTag(msg)
        case ACK:
            if ack, isIn := chans[msg.Ety.Time_stamp]; isIn{
                ack <- (msg.Ety.Msg == "agreed" )
            }
        case COMMIT:
            n.commit(msg)
        //receive message required for update
        case UPDATEINFO:
            n.log.updateLog(msg.UpdateInfo)
            fmt.Printf("receive update info from %d\n", msg.Tagval.Id)
        case UPDATEREQUEST:
            n.sendUpdate(msg.Tagval)
            fmt.Printf("sent update info to %d\n", msg.Tagval.Id)

        case IMAGE:
            fmt.Println("1")
            fmt.Printf("\tAccepted by %s\n", msg.Image)
    }
}