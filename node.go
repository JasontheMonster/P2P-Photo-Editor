package main

import (
    "fmt"
    "bufio"
    "strings"
    "os"
    "time"
)

type Node struct {
    ID          int
    addr        string
    heartbeat   int
    log         Log
    tag         Tag
    mem_list    map[int]MemListEntry
    voted       bool
    holdBack    HoldBackEty
}

// Send messages to everyone in the group
func (n *Node) broadcast(msg Message) {
    for _,entry := range n.mem_list {
        if (entry.Tag.ID != n.ID) && (entry.Active) {
            send(entry.Addr, msg)   
        }
    }
}

// send invitation to new peer (string destination address)
func (n *Node) invite(dest string) {
    fmt.Printf("\tinviting %s\n", dest)
    inv := n.createMessage(INVITE, "", n.mem_list)

    //start listening threads
    go n.ImageTransferListener()

    send(dest, inv)
}

// Peer to the network
func (n *Node) joinGroup(mem_list map[int]MemListEntry, targetId int){
    n.checkPeers(mem_list)
    tmp := map[int]MemListEntry{n.ID: n.mem_list[n.ID]}
    //ask for image
    connect_receive_image(mem_list[targetId].Addr)
    //get the image
    msg := n.createMessage(ACCEPT, "", tmp)
    n.broadcast(msg)
}

// send updated logs file list to the node who requested
func(n *Node) sendUpdate(tag Tag){
    updateLog := n.log.Entries[tag.Time_stamp:]
    msg := n.createMessageWithLog(UPDATEINFO, "", updateLog)
    send(n.mem_list[tag.ID].Addr, msg)
}

// broadcast msg to all and wait for ack
func (n *Node) updateToAll(msg Message, ack chan bool){
    n.broadcast(msg)

    acks := 1
    decs := 0
    if n.voted {
        acks = 0
        decs = 1
        n.voted = true
    }
        
    for (acks < len(n.mem_list)/2 + 1) && (decs < len(n.mem_list)/2 + 1){
        if <- ack {
            acks += 1
        } else{
            decs += 1
        }
    }

    if acks >= len(n.mem_list)/2 + 1 {
        fmt.Printf("Commited: %s\n", msg.Ety.Msg)
        n.tag.Time_stamp += 1
        n.log.append(msg.Ety)
        go sendToFront(msg.Ety.Msg)
        commit := n.createMessage(COMMIT, msg.Ety.Msg, make(map[int]MemListEntry))
        n.broadcast(commit)
    }
    delete(chans, msg.Tag.Time_stamp)
}

// broadcast heartbeat to all
func (n *Node) sendHeartbeat(done chan bool) {
    var now int64
    for {
        time.Sleep(10000 * time.Millisecond)

        n.heartbeat++
        entry := n.mem_list[n.ID]
        entry.Heartbeat = n.heartbeat
        n.mem_list[n.ID] = entry
        msg := n.createMessage(HEARTBEAT, "HB", n.mem_list)
        n.broadcast(msg)

        now = time.Now().UnixNano()
        if n.voted && (now - n.holdBack.Time > TFAIL) {
            n.voted = false
        }

        for id, entry := range n.mem_list {
            if (now - entry.Timestamp > TCLEANUP) {
                n.delNode(id)
            } else if (now - entry.Timestamp > TFAIL) {
                n.mem_list[id].deactiveNode()
            }
        }
    }
    done <- true
}

// take user input and make action
func (n *Node) userInput(done chan bool) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter Commands: ")
    for {
        text,_ := reader.ReadString('\n')
        text = strings.Replace(text, "\n", "", -1)
        input := strings.SplitN(text, " ", 2)
        switch input[0] {
            case "invite":
                go n.invite(input[1])
            case "send":
                msg := n.createDataMessage(PUBLIC, input[1])
                chans[msg.Ety.Time_stamp] = make(chan bool)
                go n.updateToAll(msg, chans[msg.Ety.Time_stamp])
        }
    }
    done <- true
}