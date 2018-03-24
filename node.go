package main

import (
    "encoding/json"
    "fmt"
    "net"
    "sync"
)

var (
    mutex = new(sync.Mutex) //lock

)

type Node struct {
    ID          int
    addr        string
    mem_list    []Node
    active_mem  map[int]bool
    conn_list   map[int]net.Conn
}

//Invitation message 
type Invite struct{
    accept     bool
    mem_list   []Node  
}

type Message struct {
    helloworld    string
    QUIT          bool   
}



func (node *Node) delNode(id int) {
    node.conn_list[id].Close()
    node.active_mem[id] = false
}


//add Peer to the network
func (node *Node) joinGroup(inv Invite){
    for _, peer := range inv.mem_list{
        //connect to the peer
        conn := connectPeer(peer.addr)
        //append peer node to the mem_list
        node.mem_list = append(node.mem_list, peer)
        //set active mem map with id = true
        node.active_mem[peer.ID] = true
        //store connection in a map
        node.conn_list[peer.ID] = conn
    }
}

func (node *Node) isAlive(id int) bool {
    status, prs := node.active_mem[id]
    return status || prs
}

//receive server
func (node *Node) receive(){
    tcpAddr, err := net.ResolveTCPAddr("tcp4", ":6666")
    if err != nil{
        fmt.Println("Err getting TCP.")
    }
    
    listener, err2 := net.ListenTCP("tcp", tcpAddr)
    if err2 != nil{
        fmt.Println("Err start listening.")
    }
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }

        go handleMsg(conn)
    }  
}

//function to handle message
func handleMsg(conn net.Conn){
    dec := json.NewDecoder(conn)
    msg:= new(Message)
    defer conn.Close()
    for {
        err := dec.Decode(&msg)
        if msg.QUIT {
            return
        }else{
        fmt.Println(msg.helloworld, err)
        }
    }
}

// create a connection to peer (string destination address)
func connectPeer(addr string) net.Conn{
    tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil{
        fmt.Println("Err getting addr.")
    }
    
    conn, err2 := net.DialTCP("tcp", nil, tcpAddr)
    if err2 != nil{
        fmt.Println("Err connection addr.")
    }
  
    return conn 
}

// send invitation to new peer (string destination address)
func (node *Node) invite(dest string) {
    inv := Invite{accept: false, mem_list: node.mem_list}
    conn := connectPeer(dest)
    enc := json.NewEncoder(conn)
    enc.Encode(inv)
}

// Send messages to everyone in the group
func (node *Node) sendToAll(msg Message) {
    for _,conn := range node.conn_list {
        enc := json.NewEncoder(conn)
        enc.Encode(msg)
    }
}

func main() {
    var (
        node Node
        memlist []Node
    )
    node.ID = 1
    node.addr = "127.0.0.1:8888"
    node.mem_list = memlist
    node.active_mem = make(map[int]bool)
    node.conn_list = make(map[int]net.Conn)
}
















