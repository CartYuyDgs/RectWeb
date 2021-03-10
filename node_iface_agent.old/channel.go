package main

import (
	"sync"

	"mecas/pkg/errors"

	"github.com/gorilla/websocket"
)

var errClosed = errors.New("conn is closed")

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan [0]byte
	mutex     sync.Mutex
	isClosed  bool
}

func CreateConnection(wsConn *websocket.Conn) (conn *Connection) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 128),
		outChan:   make(chan []byte, 128),
		closeChan: make(chan [0]byte, 1),
	}

	go conn.readLoop()
	go conn.writeLoop()

	return
}

func (conn *Connection) ReadMessage() (msg []byte, err error) {
	select {
	case msg = <-conn.inChan:
	case <-conn.closeChan:
		err = errClosed
	}
	return
}

func (conn *Connection) WriteMessage(msg []byte) (err error) {
	select {
	case conn.outChan <- msg:
	case <-conn.closeChan:
		err = errClosed
	}
	return
}

func (conn *Connection) Close() {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if !conn.isClosed {
		_ = conn.wsConn.Close()
		close(conn.closeChan)
		conn.isClosed = true
	}
}

func (conn *Connection) readLoop() {
	var (
		msg []byte
		err error
	)

outter:
	for {
		if _, msg, err = conn.wsConn.ReadMessage(); err != nil {
			break outter
		}
		select {
		case conn.inChan <- msg:
		case <-conn.closeChan:
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
		case msg = <-conn.outChan:
			if err = conn.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break outter
			}
		case <-conn.closeChan:
			break outter
		}
	}

	conn.Close()
}
