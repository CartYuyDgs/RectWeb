package main

import "sync"

type hub struct {
	sessions   map[*Connection]bool
	//broadcast  chan *envelope
	register   chan *Connection
	unregister chan *Connection
	//exit       chan *envelope
	open       bool
	rwmutex    *sync.RWMutex
}

