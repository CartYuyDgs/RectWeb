package main

import (
	"context"
	"fmt"
	"time"
)

type One struct {
	name string
	cancel context.CancelFunc
	ctx context.Context
	status bool
}

type HostMsg struct {
	hostname string
	chanMsg chan *One
}

func CreatHost() *HostMsg {
	host := &HostMsg{
		hostname: "abc",
		chanMsg:  make(chan *One),
	}
	return host
}

func CreateCTX() (context.Context, context.CancelFunc) {

	timeout := 4 * time.Second
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	return ctx, cancle
}

func (o *One)Start() {
	fmt.Println(o.name)
	fmt.Println(o.status)

	go func() {
		for {
			fmt.Println("aaaaaaaaaaaaaaaaaaa")
			time.Sleep(time.Second)
		}

	}()

loop:
	for {
		select {
		case <-o.ctx.Done():
			fmt.Println("over ",o.name)
			break loop
		}
	}
}

func CreateOne() *One {
	ctx, cancle := CreateCTX()

	one :=&One{
		name:   "ONE",
		cancel: cancle,
		ctx:    ctx,
		status: true,
	}

	return one
}