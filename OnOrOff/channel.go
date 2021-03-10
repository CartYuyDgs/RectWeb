package main

func (h *HostMsg)Loop() {

	for {
		select {
		case one := <-h.chanMsg:
			one.Start()
		}
	}
}
