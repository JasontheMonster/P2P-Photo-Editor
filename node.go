package main

import (
    "fmt"
    "bufio"
    "strings"
    "os"
    "time"
    "encoding/base64"
)

type Node struct {
    ID          int
    addr        string
    log         Log
    tag         Tag
    mem_list    map[int]string
    active_mem  map[int]bool
}

func (n *Node) deactiveNode(id int) {
    n.active_mem[id] = false
}

func (n *Node) delNode(id int) {
    delete(n.mem_list, id)
    delete(n.active_mem, id)
}

func (n *Node) isAlive(id int) bool {
    status, prs := n.active_mem[id]
    return status || prs
}

func (n *Node) checkPeers(memlist map[int]string) {
    for id, addr := range memlist{
        if _, isIn := n.mem_list[id]; !isIn {
            n.mem_list[id] = addr
            n.active_mem[id] = true
        }
    }
    return
}

//add Peer to the network
func (n *Node) joinGroup(mem_list map[int]string){
	for id, addr := range mem_list{
		n.mem_list[id] = addr
		n.active_mem[id] = true
	}
    //n.checkPeers(mem_list)
    arg := "Invitation accepted by: " + n.addr
    msg := n.createMessage(ACCEPT, arg, n.mem_list)
    n.broadcast(msg)
}

// send invitation to new peer (string destination address)
func (n *Node) invite(dest string) {
    fmt.Printf("\tinviting %s\n", dest)
    inv := n.createMessage(INVITE, "invite", n.mem_list)
    //inv.Image = n.encodeImage()
    //fmt.Println(inv.Image)
    send(dest, inv)
    //fmt.Println("1")
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
    for _,addr := range n.mem_list {
	   send(addr, msg)
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
                msg := n.createMessage(PUBLIC, input[1], n.mem_list)
                go n.broadcast(msg)
        }
    }
    done <- true
}

func (n *Node) heartbeat(done chan bool) {
    for {
        time.Sleep(1000 * time.Millisecond)
        msg := n.createMessage(HEARTBEAT, "HB", n.mem_list)
        n.broadcast(msg)
    }
    done <- true
}