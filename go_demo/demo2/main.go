package main

import (
	"errors"
	"fmt"
	"log"
)

func main() {
	err := errors.New("eeeee")
	fmt.Printf("%+v", err)
	log.Printf("%+v", err)

	fmt.Printf("%+v", "hahahahahahahh")
	log.Printf("%+v", "hahahahahahh")
}
