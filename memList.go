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

func (m MemListEntry) deactiveNode() {
    m.Active = false
}

func (n *Node) delNode(id int) {
    delete(n.mem_list, id)
}

func (n *Node) isAlive(id int) bool {
    ety, prs := n.mem_list[id]
    if prs {
        return true
    }
    return ety.Active
}

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