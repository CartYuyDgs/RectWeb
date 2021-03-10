package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func CreateCTX() (context.Context,context.CancelFunc) {

	timeout := 5 * time.Second
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	return ctx,cancle
}

func GoToSm() {
	fmt.Println("to do something....")
}

type Abc struct {
	name string
	cancle context.CancelFunc
	c chan int
	ctx context.Context
}

func (a *Abc)toStart() {
	//端口UP,down
	fmt.Println("start...........")
	fmt.Println("up...........")
	//读状态，每1秒
	time.Sleep(time.Second)
	a.ctx, a.cancle = CreateCTX()
	GoToSm()

	for {
		select {
		case <- a.ctx.Done():
			fmt.Println("over ........")
			fmt.Println("down ..........")
			a.c <-0
			return
		default:
			time.Sleep(time.Second * 1)
		}
	}


}

func (a *Abc)toStop() {
	a.cancle()
	//端口down
	fmt.Println("down ..........")
	fmt.Println("stop ..........")

}

func main() {

	a := Abc{}
	a.c = make(chan int, 2)
	fmt.Println("ssssmain.......")

	go func(a *Abc) {
		for {
			//log.Println("for .......")
			select {
			case value, ok := <-a.c:
				log.Println("select.......", ok)
				if value == 1 {
					a.toStart()
				} else if value == 2{
					fmt.Println("to stop.......  ")
					a.toStop()
				} else {
					fmt.Println("stop ")
					break
				}
			default:
				//log.Println("defalut.......")
				time.Sleep(1 * time.Second)
			}
		}
	}(&a)

	a.c <- 1
	fmt.Println("dddddd.......")
	time.Sleep(time.Second * 8)
	//a.c <- 2

	time.Sleep(time.Second*20)


}
