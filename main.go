package main

import (
	"sync"
	"flag"
    "time"
)

var (
    mutex = new(sync.Mutex) //lock
    chans = make(map[int](chan bool))
)

func main() {
    done := make(chan bool)
    var node Node
    node.heartbeat = 0;
    //command line input node id & address
    flag.IntVar(&node.ID, "id", 0, "specify the node id")
    flag.StringVar(&node.addr, "addr", "127.0.0.1:8080", "specify the node address")
    flag.Parse()
    //initialized tag of the node {id, timestamp=0}
    node.tag = createTag(node.ID, 0)
    //initilaize first log with empty entry and timestamp=0
    node.log = initLog(0)
    node.voted = false
    //initialize membership list and connection list
    node.mem_list = make(map[int]MemListEntry) 
    //put itself in the list
    entry := MemListEntry{ID: node.ID, Addr: node.addr, Heartbeat: node.heartbeat, Tag: node.tag, Timestamp: time.Now().UnixNano(), Active: true}
    node.mem_list[node.ID] = entry

    //listening thread
    go node.server(done)
    //user input thread
    go node.userInput(done)
    go node.sendHeartbeat(done)
    for i := 0; i < 3; i++{
        <- done
    }
}
