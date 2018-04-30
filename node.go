package main

import (
    "fmt"
    "bufio"
    "strings"
    "os"
    "time"
)

type Node struct {
    ID          int                     // id
    addr        string                  // ip address
    heartbeat   int                     // number of heartbeat
    log         Log                     // log entries
    tag         Tag                     // tag value
    mem_list    map[int]MemListEntry    // membership list
    voted       bool                    // if current node has voted for a public message
    holdBack    int64             // hold back time for pre-commit phase
    localsendAddr   string              // tcp address to send to the front end
    localrecAddr    string              // tcp address to receive from the front end
    Image_path  string                  // local path to the chosen image
    HasImage    bool                    // if current node has w image working in progress
}

// Send messages to every active node in the group except self (message to send)
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
    // send the invite message
    send(dest, inv)
}

// broadcast join message to all active nodes and receive image from inviter (membership list, inviter id)
func (n *Node) joinGroup(mem_list map[int]MemListEntry, targetId int){
    n.checkPeers(mem_list)
    tmp := map[int]MemListEntry{n.ID: n.mem_list[n.ID]}
    // receive the image from inviter
    n.connect_receive_image(mem_list[targetId].Addr)
    // broadcast join group message
    msg := n.createMessage(ACCEPT, "", tmp)
    n.broadcast(msg)
}

// send updated logs file list to the node who requested (local tag)
func(n *Node) sendUpdate(tag Tag){
    updateLog := n.log.Entries[tag.Time_stamp:]
    msg := n.createMessageWithLog(UPDATEINFO, "", updateLog)
    send(n.mem_list[tag.ID].Addr, msg)
}

// broadcast msg to all and wait for ack for majority of active nodes
func (n *Node) updateToAll(msg Message, ack chan bool){
    n.broadcast(msg)

    // init number of votes and declines
    acks := 1
    decs := 0
    if n.voted { // if self already voted
        acks = 0
        decs = 1
    }
    // self vote
    n.voted = true
    
    // get quorum size by majority of active nodes
    quorumSize := 0
    for _,mem := range n.mem_list {
        if mem.Active {
            quorumSize += 1
        }
    }
    quorumSize = (quorumSize + 1) / 2
    
    // wait for quorumsize of nodes to reply ack or enough declines to make it impossible
    for (acks < quorumSize) && (decs < len(n.mem_list)-quorumSize){
        select {
        case res := <-ack:
            if res {
                acks += 1
            } else {
                decs += 1
            }
        case <-time.After(3 * time.Second): // if no receiving response in 4 seconds, break the loop and abort this update
            fmt.Println("timeout")
            break
        }
    }

    // if more than or equal to quorumsize of nodes replied ack
    if acks >= quorumSize {
        // self commit
        fmt.Printf("Commited by self: %s, %d\n", msg.Ety.Msg, msg.Ety.Time_stamp)
        n.tag.Time_stamp += 1
        n.log.append(msg.Ety)
        n.applyLog()
        n.voted = false
        // broadcast commit message
        commit := n.createMessage(COMMIT, msg.Ety.Msg, make(map[int]MemListEntry))
        n.broadcast(commit)
    } else { // abort this update
        n.voted = false
    }
    // delete the channel for this update
    delete(chans, msg.Tag.Time_stamp)
}

// broadcast heartbeat to all
func (n *Node) sendHeartbeat(done chan bool) {
    var now int64
    for {
        // wait 10 second
        time.Sleep(10000 * time.Millisecond)

        // broadcasr heartbeat message
        n.heartbeat++
        entry := n.mem_list[n.ID]
        entry.Heartbeat = n.heartbeat
        n.mem_list[n.ID] = entry
        msg := n.createMessage(HEARTBEAT, "HB", n.mem_list)
        n.broadcast(msg)

        // if not reciving commit message for TFAIL time, abort pending voted message
        now = time.Now().UnixNano()
        if n.voted && (now - n.holdBack > TFAIL) {
            n.voted = false
        }

        // for all nodes in membership list
        for id, entry := range n.mem_list {
            if (id != n.ID) {
                if (now - entry.Timestamp > TCLEANUP) { // if not reciving heartbeat from certain node for TCLEANUP time, delete it from list
                    n.delNode(id)
                } else if (now - entry.Timestamp > TFAIL) { // if not reciving heartbeat from certain node for TFAIL time, mark it not active
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