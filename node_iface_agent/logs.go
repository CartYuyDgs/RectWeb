package main

import (
	"log"
	"os"
)

var loginfo *log.Logger
var logWarning *log.Logger
var logError *log.Logger


func init() {
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}

	loginfo = log.New(file, "[host log info:]",log.LstdFlags | log.Lshortfile | log.LUTC)
	logWarning = log.New(file, "[host log Warning:]",log.LstdFlags | log.Lshortfile | log.LUTC)
	logError = log.New(file, "[host log Error:]",log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}




