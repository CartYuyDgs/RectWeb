package melody

import "sync"

type hub struct {
	sessions 	map[*Session]bool
	register 	chan *Session
	unregister 	chan *Session
	exit 		chan *Session
	open 		bool
	rwmutex 	*sync.RWMutex
}


func NewHub() *hub {
	return &hub{
		sessions: make(map[*Session]bool),
		register: make(chan *Session),
		unregister: make(chan *Session),
		open: true,
		rwmutex: &sync.RWMutex{},
	}
}


func (h *hub)Run() {

loop:
	for  {
		select {
		case m := <-h.register:
			h.rwmutex.Lock()
			h.sessions[m] = true
			h.rwmutex.Unlock()
		case m := <- h.unregister:
			h.rwmutex.Lock()
			delete(h.sessions, m)
			h.rwmutex.Unlock()
		case  <- h.exit:
			h.rwmutex.Lock()
			for s:= range h.sessions {
				delete(h.sessions, s)
				//s.Close()
			}
			h.open = false
			h.rwmutex.Unlock()
			break loop
		}
	}
}
