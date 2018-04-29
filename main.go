package main

import (
	"sync"
	"flag"
    "time"
    // "fmt"
)

var (
    // mutex lock
    mutex = new(sync.Mutex)
    // map of chanels for send process
    chans = make(map[int](chan bool))
)

func main() {
    // chanel to mark completion of the process
    done := make(chan bool)
    // init node object
    var node Node
    // set number of heartbeat to 0
    node.heartbeat = 0;
    // command line input node id
    flag.IntVar(&node.ID, "id", 0, "specify the node id")
    // command line input node addr
    flag.StringVar(&node.addr, "addr", "127.0.0.1:8080", "specify the node address")
    flag.StringVar(&node.localrecAddr, "listenAddr", "127.0.0.1:5050", "listening address")
    flag.StringVar(&node.localsendAddr, "sendAddr", "127.0.0.1:5051", "sending address")
    // parse command line input
    flag.Parse()
    // initialized tag of the node {id, timestamp=0}
    node.tag = createTag(node.ID, 0)
    // initilaize first log with empty entry and timestamp=0
    node.log = initLog(0)
    // init voted to false -> not voted yet
    node.voted = false
    // initialize membership list and connection list
    node.mem_list = make(map[int]MemListEntry) 
    // put itself in the map
    node.mem_list[node.ID] = MemListEntry{Addr: node.addr, Heartbeat: node.heartbeat, Tag: node.tag, Timestamp: time.Now().UnixNano(), Active: true}

    //fmt.Println("initialize node", node.mem_list)
    node.Image_path = ""
    go node.localConnection(node.localrecAddr)
    //listening thread
    go node.server(done)
    //user input thread
    // go node.userInput(done)
    // heartbeat thread
    go node.sendHeartbeat(done)
    // wait for threads to finish
    for i := 0; i < 2; i++{
        <- done
    }
}
