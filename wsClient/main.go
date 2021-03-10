package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/CartYuyDgs/RectWeb/wsClient/mode"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

const wsUrl = "127.0.0.1:8888"
const WsPath = "/ws"

func main() {

	u := url.URL{
		Scheme: "ws",
		Host:   wsUrl,
		Path:   WsPath,
	}

	hostname := "mec54"

	header := make(http.Header, 4)
	header.Set("WebSocketAgent", hostname)
	var username, password string
	username = "bonc"
	password = "bonc@123"
	header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
	//WsConn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	//WsConn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	//WsConn, _, err :=websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}.Dial(u.String(), header)
	log.Println("hostname: ", hostname, " ", "", wsUrl)
	dialer := websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}}
	WsConn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		panic(err)
	}

	defer WsConn.Close()

	conn := CreateConnection(WsConn)
	conn.isClosed = false

	go receiveCommand(conn)

	for {
		msgdata, err := getMessage()
		if err != nil {
			log.Fatalln(err)
		}

		err = sendMessage(conn, msgdata)
		if err != nil {
			log.Fatalln(err)
		}

		time.Sleep(4 * time.Second)
	}

}

func getMessage() (MsgDate []byte, err error) {
	msg := mode.NodeMessage{
		NodeName:         "m2",
		CPUavailability:  21.9,
		MemAvailability:  12.1,
		DiskAvailability: 22.3,
		NumNetwork:       2,
		NodeStatus:       true,
		NodeType:         mode.NodeAvailable,
		OwnCluser:        "abc",
	}

	MsgDate, err = json.Marshal(msg)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func sendMessage(conn *Connection, MsgDate []byte) error {
	err := conn.WriteMessage(MsgDate)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

// 接收服务端发回的指令，并执行指令
func receiveCommand(conn *Connection) {
	var (
		err error
		msg []byte
	)
	for {
		if msg, err = conn.ReadMessage(); err != nil {
			return
		}
		if err != nil {
			fmt.Println("read:", err)
			continue
		}

		fmt.Println("recv... ",string(msg))


	}
}
