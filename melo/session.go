package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Session struct {
	Request *http.Request
	keys map[string]interface{}
	conn *websocket.Conn
	melody  *Meloy
	open bool
	rwmutex	*sync.RWMutex
}


func (s *Session)writeMessage(message *envelope) {

}

func (s *Session) closed() bool {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()

	return !s.open
}

func (s *Session) close() {
	if !s.closed() {
		s.rwmutex.Lock()
		s.open = false
		s.conn.Close()
		//close(s.output)
		s.rwmutex.Unlock()
	}
}


func (s *Session) Close() error {
	if s.closed() {
		return errors.New("session is already closed")
	}

	s.writeMessage(&envelope{t: websocket.CloseMessage, msg: []byte{}})

	return nil
}
