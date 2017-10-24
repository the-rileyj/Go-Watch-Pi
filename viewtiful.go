package main

import (
	"strings"
	"time"
	"fmt"
	"os"
	"mime/multipart"
	"net/http"
	"io"
	"bytes"
	"path/filepath"

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

func checkCamera(c chan bool) {
	out, _ := exec.Command("vcgencmd", "get_camera").Output()
	for !strings.Contains(out, "supported=1 detected=1") {
		out, _ = exec.Command("vcgencmd", "get_camera").Output()
		time.Sleep(time.Second)
	}
	c <- true
}

func checkInternet(c chan bool) {
	out, _ := exec.Command("vcgencmd", "get_camera").Output()
	for strings.Contains(out, "connect: Network is unreachable") {
		out, _ = exec.Command("vcgencmd", "get_camera").Output()
		time.Sleep(time.Second)
	}
	c <- true
}

type message struct {
	Pi      bool   `json:"pi"`
	Message string `json:"message"`
}

func main() {
	onCam, onInt := false, false
	oC, oI := make(chan bool, chan bool)
	go checkCamera(oC)
	go checkInternet(oI)

	for !(onCam && onInt) {
		select{
		case v := <- oC:
			onCam = true
		case v := <- oC:
			onInt = true	
		}
	}
	
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Problem with getting working directory", err)
	}

	client := &http.Client{}
	wsite := "ws://www.therileyjohnson.com/wsspy"
	psite := "http://www.therileyjohnson.com/subphoto"
	var dialer websocket.Dialer
	var m message
	c := camera.New("pics/")

	for {
		conn, _, err := dialer.Dial(wsite, nil)
		if err != nil {
			fmt.Println("Error astablishing websocket connection with site", err)
			time.Sleep(time.Second * 5)
			continue
		}

		for {
			if err = conn.ReadJSON(&m); err != nil {
				fmt.Println("JSON Read Error", err.Error())
				break
			}

			s, err := c.Capture()
			if err != nil {
				fmt.Printf("Problem with camera\n%s\n", err)
				break
			}
			
			req, err := newfileUploadRequest(psite, nil, "pimg", path + "/" + s)
			if err != nil {
				fmt.Println("Error Creating File Upload", err)
			}

			if _, err := client.Do(req); err != nil {
				fmt.Println("Error Making POST Request", err)
			}

			if err := conn.WriteJSON(message{true, ""}); err != nil {
				fmt.Println("JSON writing error", err)
			}

			go os.Remove(path + "/" + s)
		}
	}
}
