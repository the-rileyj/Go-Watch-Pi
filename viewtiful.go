package picamera

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/loranbriggs/go-camera"
)

type message struct {
	message string
}

func main() {
	site := "ws://therileyjohnson.com/spyer"
	var dialer websocket.Dialer
	var m message

	c := camera.New("pics/")
	conn, _, err := dialer.Dial(site, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn.ReadJSON(&m)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		s, err := c.Capture()
		if err != nil {
			fmt.Printf("Problem with camera\n%s", err)
		}
		file, _ := os.Open(s)
		defer file.Close()
		fileInfo, _ := file.Stat()
		size := fileInfo.Size()
		sbytes := make([]byte, size)
		buffer := bufio.NewReader(file)
		_, err = buffer.Read(sbytes)
		conn.WriteMessage(websocket.TextMessage, sbytes)
		os.Remove(s)
	}
}
