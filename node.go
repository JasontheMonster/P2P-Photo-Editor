package main

import (
    "encoding/json"
    "fmt"
    "net"
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
    conn_list   map[int]net.Conn
}

func (n *Node) deactiveNode(id int) {
    n.conn_list[id].Close()
    n.active_mem[id] = false
}

func (n *Node) delNode(id int) {
    delete(n.conn_list, id)
    delete(n.mem_list, id)
    delete(n.active_mem, id)
}

func (n *Node) isAlive(id int) bool {
    status, prs := n.active_mem[id]
    return status || prs
}

//creates server
func (n *Node) server(done chan bool){
    tcpAddr, err := net.ResolveTCPAddr("tcp4", n.addr)
    if err != nil{
        fmt.Println(err)
    }
    listener, err2 := net.ListenTCP("tcp", tcpAddr)
    if err2 != nil{
        fmt.Println(err2)
    }
    for {
        conn, err3 := listener.Accept()
        if err3 != nil {
		  fmt.Println(err3)
        }
        go n.handleMsg(conn)
    } 
    done <- true
}

//function to handle message
func (n *Node) handleMsg(conn net.Conn){
	var msg Message
    dec := json.NewDecoder(conn)
    defer conn.Close()
    for {
        if err := dec.Decode(&msg); err != nil {
            fmt.Println(err)
        }
		fmt.Println(msg)
        switch msg.Kind {
            case INVITE:
                n.tag.time_stamp = msg.Tagval.time_stamp
                n.log = initLog(msg.Tagval.time_stamp)
                n.joinGroup(msg.Mem_list)
            case PUBLIC:
                n.checkPeers(msg.Mem_list)
                //n.updateTag(msg.Tagval)
                fmt.Println(msg.Ety.Msg)
            case HEARTBEAT:
                n.checkPeers(msg.Mem_list)
                n.updateTag(msg.Tagval)
            case ACCEPT:
                continue
            case DECLINE:
                n.updateTag(msg.Tagval)
        }
    }
}

func (n *Node) checkPeers(memlist map[int]string) {
    for id, addr := range memlist{
        if _, isIn := n.mem_list[id]; !isIn {
            n.mem_list[id] = addr
            n.active_mem[id] = true
            n.conn_list[id] = n.connectPeer(addr)
        }
    }
    return
}

//add Peer to the network
func (n *Node) joinGroup(mem_list map[int]string){
	for id, addr := range mem_list{
		fmt.Println(addr)
		n.mem_list[id] = addr
		n.active_mem[id] = true
		n.conn_list[id] = n.connectPeer(addr)
	}
    //n.checkPeers(mem_list)
    arg := "Invitation accepted by: " + n.addr
    msg := n.createMessage(PUBLIC, arg, n.mem_list) 
    n.sendToAll(msg)
}

// create a connection to peer (string destination address)
func (node *Node) connectPeer(addr string) net.Conn{
    tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil{
        fmt.Println(err)
    }   
    conn, err2 := net.DialTCP("tcp", nil, tcpAddr)
    if err2 != nil{
        fmt.Println(err2)
    }
    return conn 
}

// send invitation to new peer (string destination address)
func (n *Node) invite(dest string) {
    inv := n.createMessage(INVITE, "", n.mem_list)
	conn := n.connectPeer(dest)
	json.NewEncoder(conn).Encode(inv)
}

// Send messages to everyone in the group
func (n *Node) sendToAll(msg Message) {
    for _,conn := range n.conn_list {
	json.NewEncoder(conn).Encode(msg)
    }
}

func (n *Node) userInput(done chan bool) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter command: ")
    for {
        text,_ := reader.ReadString('\n')
        text = strings.Replace(text, "\n", "", -1)
        input := strings.SplitN(text, " ", 2)
        switch input[0] {
            case "invite":
                go n.invite(input[1])
            case "send":
                msg := n.createMessage(PUBLIC, input[1], n.mem_list)
                go n.sendToAll(msg)
        }
    }
    done <- true
}

func (n *Node) heartbeat(done chan bool) {
    for {
        time.Sleep(1000 * time.Millisecond)
        msg := n.createMessage(HEARTBEAT, "HB", n.mem_list)
        n.sendToAll(msg)
    }
    done <- true
}
