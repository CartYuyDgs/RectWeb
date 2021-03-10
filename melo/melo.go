package main

import (
	"github.com/gorilla/websocket"
)

type handlerMessageFunc func(*Session, []byte)
type handleCloseFunc func(*Session, int, string) error
type handleSessionFunc func(*Session)
type filterFunc func(*Session) bool



type Meloy struct {
	Upgrader           *websocket.Upgrader
	messageHandler     handlerMessageFunc
	messageSendHandler handlerMessageFunc
	connectHandler     handleSessionFunc
	disconnectHandler  handleSessionFunc
	hub                *hub
}
