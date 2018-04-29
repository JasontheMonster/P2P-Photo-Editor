package main 

import (
    "net"
    "fmt"
    "strings"
    "os"
)


func (n *Node) localConnection(addr string){

	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
    if err != nil{
        fmt.Println("User not found")
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
        fmt.Println(s)
        //fmt.Println("receive front end invite and send", n.mem_list)
        n.invite(s)
    } else if (strings.HasPrefix(s, "PATH")){
        s = strings.TrimPrefix(s, "PATH:")
        fmt.Println("this is the image path", s)
        n.Image_path = s
        n.HasImage = true
    } else if (strings.HasPrefix(s, "quit")) {
        n.sendToFront("quit")
        n.HasImage = false
        os.Exit(0)
    } else {
        msg := n.createDataMessage(PUBLIC, s)
        chans[msg.Ety.Time_stamp] = make(chan bool)
        n.updateToAll(msg, chans[msg.Ety.Time_stamp])
    }
}

func (n *Node) sendToFront(logEty string) {
    conn,_ := net.Dial("tcp", n.localsendAddr)
    defer conn.Close()
    conn.Write([]byte(logEty))
}