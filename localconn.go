package main 

import (
    "net"
    "fmt"
    "strings"
)


func (n *Node) localConnection(addr string){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
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
        go n.handleClient(conn)

    } 
    //done <- true
}

func (n *Node) handleClient(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 1024)
    relen, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }

    s := string(buf[:relen])
    if (strings.HasPrefix(s,"invite")) {
        s = strings.TrimPrefix(s,"invite")
        n.invite(s)
    } else{
        msg := n.createDataMessage(PUBLIC, s)
        chans[msg.Ety.Time_stamp] = make(chan bool)
        n.updateToAll(msg, chans[msg.Ety.Time_stamp])
    }
}

func sendToFront(logEty string) {
    conn,_ := net.Dial("tcp", "127.0.0.1:5006")
    defer conn.Close()
    conn.Write([]byte(logEty))
}