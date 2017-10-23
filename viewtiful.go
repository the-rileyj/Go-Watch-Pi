package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/loranbriggs/go-camera"
)

type message struct {
	Pi      bool   `json:"pi"`
	Message string `json:"message"`
	Pic     []byte `json:"pic"`
}

func main() {
	site := "ws://www.therileyjohnson.com/wsspy"
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
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Problem with getting working directory\n%s\n", err)
				break
			}
			file, err := os.Open(wd + "/" + s)
			if err != nil {
				fmt.Printf("Problem with file\n%s\n", err)
				break
			}
			fileInfo, _ := file.Stat()
			size := fileInfo.Size()
			sbytes := make([]byte, size)
			buffer := bufio.NewReader(file)
			_, err = buffer.Read(sbytes)
			conn.WriteMessage(websocket.TextMessage, sbytes)
			/*if err := conn.WriteJSON(message{true, "", sbytes}); err != nil {
				fmt.Println("JSON writing error", err)
			}*/
			file.Close()
			go os.Remove(s)
		}
	}
}
