package main

import (
	"fmt"
	"time"
)

func main()  {
	go printHello()
	time.Sleep(time.Millisecond)
}

func printHello() {
	fmt.Println("hello")
}