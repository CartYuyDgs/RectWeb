package main

import (
	"github.com/gorilla/websocket"
)

type handleMessageFunc func(*Connection, []byte)
type handleErrorFunc func(*Connection, error)
type handleCloseFunc func(*Connection, int, string) error
type handleSessionFunc func(*Connection)
type filterFunc func(*Connection) bool

//websocket manager
type Melody struct {
	Upgrader                 *websocket.Upgrader
	messageHandler           handleMessageFunc
	messageHandlerBinary     handleMessageFunc
	messageSentHandler       handleMessageFunc
	messageSentHandlerBinary handleMessageFunc
	errorHandler             handleErrorFunc
	closeHandler             handleCloseFunc
	connectHandler           handleSessionFunc
	disconnectHandler        handleSessionFunc
	pongHandler              handleSessionFunc
	hub                      *hub
}
