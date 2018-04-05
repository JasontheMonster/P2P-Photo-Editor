package main 

import (
	"net"
	"fmt"
	"log"
)

func send(addr string, msg Message){
	udpAddr,err1 := net.ResolveUDPAddr("udp4", addr)
	if err1 != nil {
		fmt.Println("Err getting addr")
		log.Fatal(err1)
	}

	conn,err2 := net.DialUDP("udp", nil, udpAddr)
	if err2 != nil {
		fmt.Println("Err dialing.")
		log.Fatal(err2)
	}
	defer conn.Close()

	conn.Write([]byte(b))
}

//creates listener
func (n *Node) ImageTransferListener(done chan bool){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", n.addr)
    if err != nil{
        fmt.Println(err)
    }
    listener, err2 := net.ListenTCP("tcp", tcpAddr)
    if err2 != nil{
        fmt.Println(err2)
    }
    defer listener.Close()
    for {
        conn, err3 := listener.Accept()
        if err3 != nil {
		  fmt.Println(err3)
        }
        go n.handleImage(conn)
    } 
    done <- true
}

func (n *Node) handleImage(conn net.Conn){
	fmt.Println("Waiting for a image!")
	defer connection.Close()

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	//read the size of file size
	connection.Read(bufferFileSize)
	filesize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

}

