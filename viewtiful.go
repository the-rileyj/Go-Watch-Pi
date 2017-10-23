package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
    "mime/multipart"

	"github.com/gorilla/websocket"
	"github.com/loranbriggs/go-camera"
)


//From: https://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
  
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
  
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
  
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
  }

type message struct {
	Pi      bool   `json:"pi"`
	Message string `json:"message"`
}

func main() {
	site := "ws://www.therileyjohnson.com/wsspy"
	var dialer websocket.Dialer
	var m message
	c := camera.New("pics/")
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Problem with getting working directory\n%s\n", err)
		break
	}
	for {
		conn, _, err := dialer.Dial(site, nil)
		if err != nil {
			log.Fatal(err)
		}

		for {
			var b bytes.Buffer
			w := multipart.NewWriter(&b)
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
			newfileUploadRequest("", nil, "pimg", wd + "/" + s)
			// file, err := os.Open(wd + "/" + s)
			// if err != nil {
			// 	fmt.Printf("Problem with file\n%s\n", err)
			// 	break
			// }

			// fw, err := w.CreateFormFile("image", wd + "/" + s)
			// if err != nil {
			// 	fmt.Printf("Problem with adding file to form\n%s\n", err)
			// 	return 
			// }

			// if _, err = io.Copy(fw, f); err != nil {
			// 	fmt.Printf("Problem with copying file to form\n%s\n", err)
			// 	return
			// }

			// if fw, err = w.CreateFormField("pim"); err != nil {
			// 	return
			// }
			// if _, err = fw.Write([]byte("PIM")); err != nil {
			// 	return
			// }

			// w.Close()
			// req, err := http.NewRequest("POST", url, &b)
			// if err != nil {
			// 	return 
			// }
			// fileInfo, _ := file.Stat()
			// size := fileInfo.Size()
			// sbytes := make([]byte, size)
			// buffer := bufio.NewReader(file)
			// _, err = buffer.Read(sbytes)
			//conn.WriteMessage(websocket.TextMessage, sbytes)
			/*if err := conn.WriteJSON(message{true, "", sbytes}); err != nil {
				fmt.Println("JSON writing error", err)
			}*/
			go os.Remove(wd + "/" + s)
		}
	}
}
