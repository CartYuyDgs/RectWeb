package main

import (
	"fmt"
	"time"
)

func main() {
	host := CreatHost()

	go host.Loop()

	o := CreateOne()

	host.chanMsg <- o

	//time.Sleep(3 * time.Millisecond)

	fmt.Println("line")
	//o.cancel()

	time.Sleep(5 * time.Second)
}
