package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"picamera"
)

type message struct {
	message string
}

func main() {
	var dialer websocket.Dialer
	var m message
	c := camera.New("pics")
	conn, _, err := dialer.Dial("ws://therileyjohnson.com/spyer", nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn.ReadJSON(&m)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		s, err := c.Camera.Capture()
		if err != nil {
			printf("Problem with camera\n%s", err)
		}
		file, _ := os.Open("/pics/pic.png")
		defer file.Close()
		fileInfo, _ := file.Stat()
		size := fileInfo.Size()
		sbytes := make([]byte, size)
		buffer := bufio.NewReader(file)
		_, err = buffer.Read(sbytes)
		conn.WriteMessage(websocket.TextMessage, sbytes)
	}
}
