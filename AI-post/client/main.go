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
	Token       string   `json:"token"`
	Msg         string   `json:"msg"`
	Photobase64 []string `json:"data"`
}

func main() {

	//获取参数
	//var url = flag.String("url", "http://172.16.96.57:12345/detection", "input your url")
	var url = flag.String("url", "http://172.16.96.57:12345/meter", "input your url")
	var token = flag.String("token", "bonccvlab", "input your token")
	var msg = flag.String("msg", "1.jpg", "input your photo formate")
	var photo = flag.String("photo", "./photo4.jpg", "input your photo address")
	flag.Parse()
	//fmt.Println(*token, *msg, *photo)

	//now_time := time.Now()main

	for i := 0; i < 9; i++ {

		now_time := time.Now()
		CreateInfo(*url, *token, *msg, *photo)
		fmt.Printf("执行消耗的时间为:%v", now_time.Sub(time.Now()))
	}

	return

}

func CreateInfo(url, token, msg, photo string) {
	Objfile, err := os.Open(photo)
	defer Objfile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//base64
	buffer := make([]byte, 500000)
	n, _ := Objfile.Read(buffer)
	encodestring := base64.StdEncoding.EncodeToString(buffer[:n])

	//fmt.Println(encodestring)

	var info = date{
		Msg:   msg,
		Token: token,
	}

	info.Photobase64 = append(info.Photobase64, encodestring)

	//post请求
	contentType := "application/json"
	//url := "http://127.0.0.1:8000"

	res, err := postSend(url, info, contentType)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("result: ", res)
}

func postSend(url string, info date, contentType string) (string, error) {
	client := &http.Client{}
	//fmt.Println(info)
	jsonStr, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Printf("%v",jsonStr)
	//fmt.Println(len(jsonStr))
	//fmt.Println(bytes.NewBuffer(jsonStr))
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//fmt.Println(resp.StatusCode)
	return string(result), err
}
