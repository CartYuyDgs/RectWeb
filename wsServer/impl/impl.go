package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

//主机信息
type HostInfo struct {
	HostName        string     `json: "hostname"`
	Cpu             string     `json: "cpu"`
	Mem             string     `json: "mem"`
	Disk            string     `json: "disk"`
	OperationSystem string     `json: "operation"`
	Kernel          string     `json: "kernel"`
	NicNum          int        `json: "nicnum"`
	ManageNic       string     `json: "managenic"`
	ManageIP        string     `json: "manageIP"`
	Nci             []*NciInfo `json: "ncis"`
}

//网卡信息
type NciInfo struct {
	Pci           string `json: "pci"`
	NciDriver     string `json: "ncidriver"`
	DriverVersion string `json: "driverversion"`
	NciName       string `json: "nciname"`
	Bandwidth     string `json: "bandwidth"`
	NciModel      string `json: "ncimodel"`
	LinkStatus    string `json: "linkstatus"`
	MAC           string `json: "mac"`
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
