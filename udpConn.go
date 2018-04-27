package main 

import (
	"encoding/json"
	"net"
	"fmt"
	"log"
)

// try send msg to addr
func send(addr string, msg Message){
	b,_ := json.Marshal(msg)
	udpAddr,err1 := net.ResolveUDPAddr("udp4", addr)
	if err1 != nil {
		fmt.Println("Err getting addr")
		//log.Fatal(err1)
	}

	conn,err2 := net.DialUDP("udp", nil, udpAddr)
	if err2 != nil {
		fmt.Println("Err dialing.")
		//log.Fatal(err2)
	}
	defer conn.Close()

	conn.Write([]byte(b))
}

//creates listener
func (n *Node) server(done chan bool){
	udpAddr,err1 := net.ResolveUDPAddr("udp4", n.addr)
	if err1 != nil {
		fmt.Println("address not found")
	}

	conn,err2 := net.ListenUDP("udp", udpAddr)
	if err2 != nil {
		fmt.Println("address can't listent")
	}
	defer conn.Close()

	for {
		var msg Message
		buf := make([]byte, 1024)
		num,_,err3 := conn.ReadFromUDP(buf)
		if err3 != nil {
			log.Fatal(err3)
		}

		json.Unmarshal(buf[:num], &msg)
		go n.handleMsg(msg)
	}
	done<-true
}