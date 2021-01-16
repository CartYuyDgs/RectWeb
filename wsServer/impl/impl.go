package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type NodeMessage struct {
	NodeName         string  `json:"nodename"`
	CPUavailability  float64 `json:"cpu"`
	MemAvailability  float64 `json:"mem"`
	DiskAvailability float64 `json:"disk"`
	NumNetwork       float64 `json:"network"`
	NodeStatus       bool    `json:"nodestatus"`
	NodeType         int     `json:"nodetype"`
	OwnCluser        string  `json:"cliser"`
}

const (
	NodeAvailable   = 1
	NodeUnAvailable = 2
)

type Connection struct {
	wsConnect *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte

	mutex   sync.Mutex
	isClose bool
}

func InitConnect(ws *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConnect: ws,
		inChan:    make(chan []byte, 100),
		outChan:   make(chan []byte, 100),
		closeChan: make(chan byte, 1),
	}

	go conn.readLoop()
	//go conn.writeLoop()

	return
}

func (conn *Connection) Close() {
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	if !conn.isClose {
		close(conn.closeChan)
		conn.isClose = true
	}
	conn.mutex.Unlock()
}

func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)

	for {

		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			log.Println(err)
			goto ERR
		}

		select {
		case conn.inChan <- data:

			//log.Printf("inChan data %s\n", data)
		case <-conn.closeChan:
			goto ERR
		default:

		}
	}

ERR:
	conn.Close()
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:

		//log.Println("read message :", data)
	case <-conn.closeChan:
		err = errors.New("connect is close")

	}
	return data, nil
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
		log.Println("send message :", data)
	case <-conn.closeChan:
		err = errors.New("connect is close")

	}
	return
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}
