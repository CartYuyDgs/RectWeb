package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type date struct {
	token string
	msg string
	photobase64 []string
}

func main() {

	//获取参数
	var url = flag.String("url", "http:127.0.0.1:8000", "input your url")
	var token = flag.String("token", "bonc", "input your token")
	var msg = flag.String("msg", "jpg", "input your photo formate")
	var photo = flag.String("photo", "./abc.jpg", "input your photo address")
	flag.Parse()
	fmt.Println(*token,*msg,*photo)

	Objfile, err := os.Open(*photo)
	defer Objfile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//base64
	buffer := make([]byte, 500000)
	n, _ := Objfile.Read(buffer)
	encodestring := base64.StdEncoding.EncodeToString(buffer[:n])

	fmt.Println(encodestring)

	var info  = date{
		msg: *msg,
		token: *token,
	}

	info.photobase64 = append(info.photobase64,encodestring)


	//post请求
	contentType := "application/json"
	//url := "http://127.0.0.1:8000"

	res, err := postSend(*url, info, contentType)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("result: ",res)
	return

}

func postSend(url string, info date, contentType string) (string, error)  {
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(info)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(resp.StatusCode)
	return string(result),err
}
