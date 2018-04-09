package main

import (
    "fmt"
    "bufio"
    "os"
	"encoding/base64"
)

func (n *Node) encodeImage() string {
    imgFile, err := os.Open("example.png")
    
    if err != nil {
     fmt.Println(err)
     os.Exit(1)
    }

    defer imgFile.Close()

    fInfo, _ := imgFile.Stat()
    var size int64 = fInfo.Size()
    buf := make([]byte, size)

    // read file content into buffer
    fReader := bufio.NewReader(imgFile)
    fReader.Read(buf)

    imgBase64Str := base64.StdEncoding.EncodeToString(buf)

    return imgBase64Str
}