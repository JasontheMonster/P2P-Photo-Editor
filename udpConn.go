//P2P network is build on edp connection
package main 

import (
	"encoding/json"
	"net"
	"fmt"
	"log"
)

// try send msg to addr
func send(addr string, msg Message){
	//use json to serialized data before sending
	b,_ := json.Marshal(msg)
	udpAddr,err1 := net.ResolveUDPAddr("udp4", addr)
	if err1 != nil {
		fmt.Println("Err getting addr")
		return //gracefully deal with connection error
	}

	conn,err2 := net.DialUDP("udp", nil, udpAddr)
	if err2 != nil {
		fmt.Println("Err dialing.")
		return
	}
	defer conn.Close()

	//send
	conn.Write([]byte(b))
}

// creates listener
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
		//buffer size is 1024 bytes
		buf := make([]byte, 1024)
		num,_,err3 := conn.ReadFromUDP(buf)
		if err3 != nil {
			log.Fatal(err3)
		}

		//deserialize the received data
		json.Unmarshal(buf[:num], &msg)
		// call go rountine to handle the message
		go n.handleMsg(msg)
	}
	//channel to keep track of whether the node is alive
	done<-true
}