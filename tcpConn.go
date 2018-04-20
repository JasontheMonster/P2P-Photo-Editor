package main 

import (
	"net"
	"fmt"
	"strings"
	"strconv"
	"os"
	"io"
)

const BUFFERSIZE = 1024

//connect to inviter and receive image with TCP
func connect_receive_image(addr string){
	//connect to the inviter
	connection, err := net.Dial("tcp", addr)
	if err != nil{
		panic(err)
	}
	defer connection.Close()
	fmt.Println("connected to server, start receive file name and size")
	//maximum file name 64 characters
	bufferFileName := make([]byte, 64)
	//maximum file size is 999999999 bytes 
	bufferFileSize := make([]byte, 10)

	//receive size first
	connection.Read(bufferFileSize)
	//parse it to base 10 integer, trim filled characters
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	//receive file name
	connection.Read(bufferFileName)
	//trim extra characters
	fileName := strings.Trim(string(bufferFileName), ":")

	fmt.Println("start receiving file: ", fileName, " with size: ", fileSize)
	//create a new file descripter
	newFile, err := os.Create(LOG_PATH+fileName)
	if err != nil{
		panic(err)
	}
	defer newFile.Close()

	//keep track of number bytes received
	var receivedBytes int64
	for {
		//if all message is in one buffer
		if (fileSize - receivedBytes) < BUFFERSIZE{
			//write the file
			io.CopyN(newFile, connection, (fileSize-receivedBytes))
			//read the message
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		// write the whole buffer to file 
		io.CopyN(newFile, connection, BUFFERSIZE)
		//increment number of bytes received
		receivedBytes+= BUFFERSIZE
	}

	fmt.Println("Received!")
}

//fill the string with ":"  to certain length
func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}


//creates listener
func (n *Node) ImageTransferListener(){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", n.addr)
    if err != nil{
        fmt.Println(err)
    }
    listener, err2 := net.ListenTCP("tcp", tcpAddr)
    if err2 != nil{
        fmt.Println(err2)
    }
    defer listener.Close()
    finish_image := make(chan bool)

    for {
        conn, err3 := listener.Accept()
        if err3 != nil {
		  fmt.Println(err3)
        }
        go handleImage(conn, finish_image)
        x := <- finish_image
        //fmt.Println(x)
        if x {
        	break
        }
    } 
    //done <- true
}

func handleImage(conn net.Conn, finish_image chan bool){
	fmt.Println("A client has been connect")
	defer conn.Close()

	file, err := os.Open("images/android.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)

	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	finish_image <- true 
	return
}




