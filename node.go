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
    msg := n.createMessage(PUBLIC, arg, n.mem_list)
    n.broadcast(msg)
}

// send invitation to new peer (string destination address)
func (n *Node) invite(dest string) {
    fmt.Printf("\tinviting %s\n", dest)
    inv := n.createMessage(INVITE, "invite", n.mem_list)
    send(dest, inv)
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