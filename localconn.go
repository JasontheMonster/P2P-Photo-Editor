package main 

import (
    "net"
    "fmt"
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
        go handleClient(conn)

    } 
    //done <- true
}

func handleClient(conn net.Conn){
    defer conn.Close()
    buf := make([]byte, 1024)

    relen, err := conn.Read(buf)

    s := string(buf[:relen])
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }
    fmt.Println("received %s", s)
}