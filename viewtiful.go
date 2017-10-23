package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/loranbriggs/go-camera"
)

type message struct {
	Pi      bool   `json:"pi"`
	Message string `json:"message"`
}

func main() {
	site := "ws://therileyjohnson.com/wsspy"
	var dialer websocket.Dialer
	var m message
	c := camera.New("pics/")
	for {
		conn, _, err := dialer.Dial(site, nil)
		if err != nil {
			log.Fatal(err)
		}
		for {
			err = conn.ReadJSON(&m)
			if err != nil {
				fmt.Printf("JSON Read Error\n%s\n", err.Error())
				break
			}
			s, err := c.Capture()
			if err != nil {
				fmt.Printf("Problem with camera\n%s\n", err)
				break
			}
			file, err := os.Open(s)
			if err != nil {
				fmt.Printf("Problem with file\n%s\n", err)
				break
			}
			fileInfo, _ := file.Stat()
			size := fileInfo.Size()
			sbytes := make([]byte, size)
			buffer := bufio.NewReader(file)
			_, err = buffer.Read(sbytes)
			conn.WriteJSON(message{true, base64.StdEncoding.EncodeToString(sbytes)})
			file.Close()
			os.Remove(s)
		}
	}
}
