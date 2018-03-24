package photoEditor

import (
	"fmt"
)

type MemList struct {
	mem_list	Node[]
	active_mem	map[int]bool
}

type Node struct {
	ID			int
	IP			string
	port		int
	mem_list	MemList
	connection	Connection
}

func (n *Node) delNode() {
	n.connection.Close()
}