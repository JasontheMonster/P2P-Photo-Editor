package main

import (
    "fmt"
    "bufio"
    "strings"
    "os"
    "time"
    "encoding/base64"
)

type MemListEntry struct {
    ID          int
    Addr        string
    Heartbeat   int
    Tag         Tag
    Timestamp   int64
    Active      bool
}

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

func (m MemListEntry) deactiveNode() {
    m.Active = false
}

func (n *Node) delNode(id int) {
    delete(n.mem_list, id)
}

func (n *Node) isAlive(id int) bool {
    ety, prs := n.mem_list[id]
    if prs {
        return true
    }
    return ety.Active
}

func (n *Node) checkPeers(memlist map[int]MemListEntry) {
    for id, entry := range memlist{
        if _, isIn := n.mem_list[id]; !isIn {
            entry.Timestamp = time.Now().UnixNano()
            n.mem_list[id] = entry
        } else if (n.mem_list[id].Heartbeat <= entry.Heartbeat) {
            entry.Timestamp = time.Now().UnixNano()
            n.mem_list[id] = entry
        }
    }
    return
}

//add Peer to the network
func (n *Node) joinGroup(mem_list map[int]MemListEntry){
	for id, entry := range mem_list{
        entry.Timestamp = time.Now().UnixNano()
		n.mem_list[id] = entry
	}

    var memlist = make(map[int]MemListEntry)
    memlist[n.ID] = n.mem_list[n.ID]
    msg := n.createMessage(ACCEPT, "", memlist)
    n.broadcast(msg)
}

// send invitation to new peer (string destination address)
func (n *Node) invite(dest string) {
    fmt.Printf("\tinviting %s\n", dest)
    inv := n.createMessage(INVITE, "", n.mem_list)
    send(dest, inv)
}

func (n *Node) encodeImage() string {
    imgFile, err := os.Open("example.png")
    
    if err != nil {
     fmt.Println(err)
     os.Exit(1)
    }

    defer imgFile.Close()

    fInfo, _ := imgFile.Stat()
    var size int64 = fInfo.Size()
    buf := make([]byte, size)

    // read file content into buffer
    fReader := bufio.NewReader(imgFile)
    fReader.Read(buf)

    imgBase64Str := base64.StdEncoding.EncodeToString(buf)

    return imgBase64Str
}

// Send messages to everyone in the group
func (n *Node) broadcast(msg Message) {
    for _,entry := range n.mem_list {
        if (entry.ID != n.ID) && (entry.Active) {
            send(entry.Addr, msg)   
        }
    }
}

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
                msg := n.createDataMessage(PUBLIC, input[1], n.mem_list)
                chans[msg.Ety.Time_stamp] = make(chan bool)
                go n.updateToAll(msg, chans[msg.Ety.Time_stamp])
        }
    }
    done <- true
}

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
        commit := n.createMessage(COMMIT, msg.Ety.Msg, make(map[int]MemListEntry))
        n.broadcast(commit)
    }
    delete(chans, msg.Tagval.Time_stamp)
}

func (n *Node) commit(msg Message) {
    n.voted = false
    if (n.tag.compareTo(msg.Tagval) == -1) {
        n.log.append(n.holdBack.Ety)
    }
    n.tag.Time_stamp = msg.Tagval.Time_stamp 
    fmt.Printf("Commited: %s\n", n.holdBack.Ety.Msg)
}

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

//send updated logs file list to the node who requested
func(n *Node) sendUpdate(tag Tag){
    otherTime := tag.Time_stamp
    entries_history := n.log.entries
    updateLog := entries_history[otherTime:]
    msg := n.createMessageWithLog(UPDATEINFO, "", updateLog)
    send(n.mem_list[tag.Id].Addr, msg)
}