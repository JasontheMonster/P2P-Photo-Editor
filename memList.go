package main

import (
	"time"
)

type MemListEntry struct {
    Tag         Tag
    Addr        string
    Heartbeat   int
    Timestamp   int64
    Active      bool
}

// delete node in membership list
func (m MemListEntry) deactiveNode() {
    m.Active = false
}

// delete node from membership list
func (n *Node) delNode(id int) {
    delete(n.mem_list, id)
}

// update local membership list by membership list in incoming message
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