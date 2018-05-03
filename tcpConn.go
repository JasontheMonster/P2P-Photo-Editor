//Image Transfer Channel is build on TCP
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

	//store in ~/logs/
	dir, _ := os.Getwd() 
	path := dir + LOG_PATH + modify_filename(fileName)
	
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
	n.Image_path = path //set node's Image path to path
	n.sendToFront("Image@"+path)	//front end rendering
	n.HasImage = true	//change the state of HasImage to true

}


//everytime recevive a new image, modify the filename to avoid duplicate
func modify_filename(filename string) string{
	filenameRaw := strings.Split(filename, ".")
	newFilename := filenameRaw[0] + "1" + "." + filenameRaw[1]
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
        //go rountine to send a image
        go n.handleImage(conn, finish_image)

        //if the image finish receved, stop the socket
        x := <- finish_image
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

	//get the file size
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	//get the file name
	fileName := fillString(fileInfo.Name(), 64)

	//send file name and size first
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))

	//use the send buffer of size 1024 to send image batch by batch
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