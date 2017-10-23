package main

import (
	"fmt"
	//"bufio"
	//"fmt"
	//"log"
	"os"
	"mime/multipart"
	"net/http"
	"io"
	"bytes"
	"path/filepath"

	//"github.com/gorilla/websocket"
	//"github.com/loranbriggs/go-camera"
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
			log.Fatal("Error astablishing websocket connection with site", err)
			break
		}

		for {
			err = conn.ReadJSON(&m)
			if err != nil {
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

			if err := client.Do(req); err != nil {
				fmt.Println("Error Making POST Request", err)
			}

			if err := conn.WriteJSON(message{true, ""}); err != nil {
				fmt.Println("JSON writing error", err)
			}

			go os.Remove(path + "/" + s)
		}
	}
}
