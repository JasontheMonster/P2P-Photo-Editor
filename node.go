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
    localsendAddr   string
    localrecAddr    string
    Image_path  string
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
    fmt.Println("before prepare message", n.mem_list)
    inv := n.createMessage(INVITE, "", n.mem_list)

    //start listening threads
    go n.ImageTransferListener()

    fmt.Println("after prepare message", inv)
    send(dest, inv)
}

// Peer to the network
func (n *Node) joinGroup(mem_list map[int]MemListEntry, targetId int){
    fmt.Println(mem_list)
    n.checkPeers(mem_list)
    fmt.Println(mem_list)
    tmp := map[int]MemListEntry{n.ID: n.mem_list[n.ID]}
    //ask for image
    n.connect_receive_image(mem_list[targetId].Addr)
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
    }
    n.voted = true
    
    quorumSize := len(n.mem_list)/2
    if len(n.mem_list) % 2 != 0{
        quorumSize += 1
    }
    

    for (acks < quorumSize) && (decs < quorumSize){
        select {
        case res := <-ack:
            if res {
                acks += 1
            } else {
                decs += 1
            }
        case <-time.After(3 * time.Second):
            fmt.Println("timeout")
            break
        }
    }

    if acks >= len(n.mem_list)/2 + 1 {
        fmt.Printf("Commited: %s\n", msg.Ety.Msg)
        n.tag.Time_stamp += 1
        n.log.append(msg.Ety)
        go n.sendToFront(msg.Ety.Msg)
        commit := n.createMessage(COMMIT, msg.Ety.Msg, make(map[int]MemListEntry))
        n.voted = false
        n.broadcast(commit)
    } else {
        n.voted = false
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
            if (id != n.ID) {
                if (now - entry.Timestamp > TCLEANUP) {
                    n.delNode(id)
                } else if (now - entry.Timestamp > TFAIL) {
                    n.mem_list[id].deactiveNode()
                }
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