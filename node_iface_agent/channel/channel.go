package channel

import (
	"sync"

	"mecas/pkg/errors"

	"github.com/gorilla/websocket"
)

var errClosed = errors.New("conn is closed")

type Connection struct {
	WsConn    *websocket.Conn
	InChan    chan []byte
	OutChan   chan []byte
	CloseChan chan [0]byte
	Mutex     sync.RWMutex
	IsClosed  bool
}

func CreateConnection(wsConn *websocket.Conn) (conn *Connection) {
	conn = &Connection{
		WsConn:    wsConn,
		InChan:    make(chan []byte, 128),
		OutChan:   make(chan []byte, 128),
		CloseChan: make(chan [0]byte, 1),
	}

	go conn.readLoop()
	go conn.writeLoop()

	return
}

func (conn *Connection) ReadMessage() (msg []byte, err error) {
	select {
	case msg = <-conn.InChan:
	case <-conn.CloseChan:
		err = errClosed
	}
	return
}

func (conn *Connection) WriteMessage(msg []byte) (err error) {
	select {
	case conn.OutChan <- msg:
	case <-conn.CloseChan:
		err = errClosed
	}
	return
}

func (conn *Connection) Close() {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()
	if !conn.IsClosed {
		_ = conn.WsConn.Close()
		close(conn.CloseChan)
		conn.IsClosed = true
	}
}

func (conn *Connection) readLoop() {
	var (
		msg []byte
		err error
	)

outter:
	for {
		if _, msg, err = conn.WsConn.ReadMessage(); err != nil {
			break outter
		}
		select {
		case conn.InChan <- msg:
		case <-conn.CloseChan:
			break outter
		}
	}

	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		msg []byte
		err error
	)

outter:
	for {
		select {
		case msg = <-conn.OutChan:
			if err = conn.WsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break outter
			}
		case <-conn.CloseChan:
			break outter
		}
	}

	conn.Close()
}
