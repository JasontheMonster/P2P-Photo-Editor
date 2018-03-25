package main

import (
    "encoding/json"
    "fmt"
    "net"
    "sync"
    "flag"
    "bufio"
    "strings"
    "os"
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

type Message struct {
    Kind        string  // "INVITE", PUBLIC"
    Msg         string
    Mem_list    []Node
    //QUIT          bool   
}

func (node *Node) delNode(id int) {
    node.conn_list[id].Close()
    node.active_mem[id] = false
}


//add Peer to the network
func (node *Node) joinGroup(mem_list []Node){
    for _, peer := range mem_list{
        //connect to the peer
        conn := node.connectPeer(peer.addr)
        //append peer node to the mem_list
        node.mem_list = append(node.mem_list, peer)
        //set active mem map with id = true
        node.active_mem[peer.ID] = true
        //store connection in a map
        node.conn_list[peer.ID] = conn
    }
    fmt.Println(node.mem_list)
}

func (node *Node) isAlive(id int) bool {
    status, prs := node.active_mem[id]
    return status || prs
}

//creates server
func (node *Node) server(nonstop chan bool){
    tcpAddr, err := net.ResolveTCPAddr("tcp4", node.addr)
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

        go node.handleMsg(conn)
    } 
    nonstop <- true
}

//function to handle message
func (node *Node) handleMsg(conn net.Conn){
    dec := json.NewDecoder(conn)
    msg:= new(Message)
    defer conn.Close()
    for {
        if err := dec.Decode(msg); err != nil {
            return
        }
        switch msg.Kind {
            case "INVITE":
                node.joinGroup(msg.Mem_list)
            case "PUBLIC":
                fmt.Println(msg.Msg)
        }
    }
}

// create a connection to peer (string destination address)
func (node *Node) connectPeer(addr string) net.Conn{
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
    inv := createMessage("INVITE", "", node.mem_list)
    conn := node.connectPeer(dest)
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
    nonstop := make(chan bool)
    var (
        node Node
        memlist []Node
    )
    flag.IntVar(&node.ID, "id", 0, "specify the node id")
    flag.StringVar(&node.addr, "addr", "127.0.0.1:8080", "specify the node address")
    flag.Parse()
    node.active_mem = make(map[int]bool)
    node.conn_list = make(map[int]net.Conn)
    memlist = append(memlist, node)
    node.mem_list = memlist
    go node.server(nonstop)
    go node.userInput(nonstop)
    <- nonstop
}

func createMessage(Kind string, Msg string, Mem_list []Node) Message {
    var msg Message
    msg.Kind = Kind
    msg.Msg = Msg
    msg.Mem_list = Mem_list
    return msg
}

func (node *Node) userInput(nonstop chan bool) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter command:")
    for {
        text,_ := reader.ReadString('\n')
        text = strings.Replace(text, "\n", "", -1)
        input := strings.SplitN(text, " ", 2)
        switch input[0] {
            case "invite":
                go node.invite(input[1])
            case "send":
                msg := createMessage("PUBLIC", input[1], make([]Node,0))
                go node.sendToAll(msg)
        }
    }

    nonstop <- true
}
















