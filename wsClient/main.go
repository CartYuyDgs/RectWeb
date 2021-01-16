package main

import (
	"encoding/json"
	"github.com/CartYuyDgs/RectWeb/wsClient/mode"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

const wsUrl = "ws://127.0.0.1:8888/ws"
const origin = "http://127.0.0.1:8888/"

func main() {
	ws, err := websocket.Dial(wsUrl, "", origin)
	if err != nil {
		panic(err)
	}

	defer ws.Close()

	for {
		msgdata, err := getMessage()
		if err != nil {
			log.Fatalln(err)
		}

		err = sendMessage(ws, msgdata)
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

func sendMessage(conn *websocket.Conn, MsgDate []byte) error {
	_, err := conn.Write(MsgDate)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}
