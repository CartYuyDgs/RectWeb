package main

import (
	"sync"
)

type hub struct {
	sessions 	map[*Session]bool
	register    chan *Session
	unregister  chan *Session
	exit 		chan *envelope
	open        bool
	rwmutex     *sync.RWMutex
}

func newHub() *hub {
	return &hub{
		sessions:   make(map[*Session]bool),
		register:   make(chan *Session),
		unregister: make(chan *Session),
		exit:       make(chan *envelope),
		open:       true,
		rwmutex:    &sync.RWMutex{},
	}
}

func (h *hub) run() {

loop:
	for {
		select {
		case s := <-h.register :
			h.rwmutex.Lock()
			h.sessions[s] = true
			h.rwmutex.Unlock()
		case m := <-h.exit:
			h.rwmutex.Lock()
			for s := range h.sessions {
				s.writeMessage(m)
				delete(h.sessions, s)
				s.Close()
			}
			h.open = false
			h.rwmutex.Unlock()
			break loop

		}
	}
}
