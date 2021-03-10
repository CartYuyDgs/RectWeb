package melody

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type handlerMessageFunc func(*Session, []byte)
type handlerSentHandler func(*Session, []byte)
type handleErrorFunc func(*Session, error)
type handleCloseFunc func(*Session, int, string) error
type handleSessionFunc func(*Session)
type filterFunc func(*Session) bool

type Melody struct {
	Upgrader				*websocket.Upgrader
	messageHandler			handlerMessageFunc
	messageSendHandler		handlerSentHandler
	errorHandler            handleErrorFunc
	closeHandler            handleCloseFunc
	connectHandler          handleSessionFunc
	disconnectHandler       handleSessionFunc
	pongHandler             handleSessionFunc
	hub                     *hub
}

func New() *Melody {
	upgrade := &websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin: func(r *http.Request) bool {return true},
	}

	hub := NewHub()
	hub.Run()

	return &Melody{
		Upgrader:           upgrade,
		messageHandler:     func(*Session, []byte) {},
		messageSendHandler: func(*Session, []byte) {},
		errorHandler:       func(*Session, error) {},
		closeHandler:       nil,
		connectHandler: 	func(session *Session) {},
		disconnectHandler:  func(session *Session) {},
		pongHandler:        func(session *Session) {},
		hub:                hub,
	}
}
