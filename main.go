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

    flag.IntVar(&node.ID, "id", 0, "specify the node id")
    flag.StringVar(&node.addr, "addr", "127.0.0.1:8080", "specify the node address")
    flag.Parse()
    node.tag = createTag(node.ID, 0)
    node.log = initLog(0)
    node.active_mem = make(map[int]bool)
    node.conn_list = make(map[int]net.Conn)
    node.mem_list = make(map[int]string) 
    node.mem_list[node.ID] = node.addr
    go node.server(done)
    go node.userInput(done)
    go node.heartbeat(done)
    for i := 0; i < 3; i++{
        <- done
    }
}