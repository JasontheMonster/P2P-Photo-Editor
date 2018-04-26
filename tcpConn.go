package main 

import (
	"net"
	"fmt"
	"strings"
	"strconv"
	"os"
	"io"
	"math/rand"
)

const BUFFERSIZE = 1024

//connect to inviter and receive image with TCP
func (n *Node) connect_receive_image(addr string){
	//connect to the inviter
	connection, err := net.Dial("tcp", addr)
	if err != nil{
		fmt.Println("addr missing", addr)
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
	//create tmp disk holder
	//path := "/tmp/image_data/"
	//if _, err := os.Stat("/tmp/image_data/"); os.IsNotExist(err) {
    	//os.Mkdir("/tmp/image_data/", os.ModePerm)
	//}
	dir, _ := os.Getwd() 
	path := dir + "/logs/" + modify_filename(fileName)
	newFile, err := os.Create(path)
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
	n.Image_path = path
	n.sendToFront("Image:"+path)

}
func modify_filename(filename string) string{
	filenameRaw := strings.Split(filename, ".")
	newFilename := filenameRaw[0] + strconv.Itoa(rand.Intn(10000)) + "." + filenameRaw[1]
	return newFilename

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
        go n.handleImage(conn, finish_image)
        x := <- finish_image
        //fmt.Println(x)
        if x {
        	break
        }
    } 
    //done <- true
}

func (n *Node) handleImage(conn net.Conn, finish_image chan bool){
	fmt.Println("A client has been connect")
	defer conn.Close()

	for n.Image_path == ""{

	}

	fmt.Println("This is the path I send", n.Image_path)
	file, err := os.Open(n.Image_path)
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




