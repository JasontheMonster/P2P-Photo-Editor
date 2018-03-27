package main

import (
	"sync"
	"flag"
	"net"
)

var (
    mutex = new(sync.Mutex) //lock
)

func main() {
    done := make(chan bool)
    var node Node
    //command line input node id & address
    flag.IntVar(&node.ID, "id", 0, "specify the node id")
    flag.StringVar(&node.addr, "addr", "127.0.0.1:8080", "specify the node address")
    flag.Parse()
    //initialized tag of the node {id, timestamp=0}
    node.tag = createTag(node.ID, 0)
    //initilaize first log with empty entry and timestamp=0
    node.log = initLog(0)
    //initialize membership list and connection list
    node.active_mem = make(map[int]bool)
    node.conn_list = make(map[int]net.Conn)
    node.mem_list = make(map[int]string) 
    //put itself in the list
    node.mem_list[node.ID] = node.addr

    //listening thread
    go node.server(done)
    //user input thread
    go node.userInput(done)
    go node.heartbeat(done)
    for i := 0; i < 3; i++{
        <- done
    }
}
