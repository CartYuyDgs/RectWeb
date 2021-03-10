package main

import "github.com/CartYuyDgs/RectWeb/klogs/log2"

//func init() {
//	var logErrFile = "./files/err.log"
//	var logInfoFile = "./files/info.log"
//	var logDebugFile = "./files/debug.log"
//	var logSaveDays = 3
//
//
//	// 日志初始化
//	Logger.Config(logErrFile, Logger.ErrorLevel)
//	Logger.Config(logDebugFile, Logger.DebugLevel)
//	Logger.Config(logInfoFile, Logger.InfoLevel)
//	Logger.SetSaveMode(logs.BySize)
//	Logger.SetSaveDays(logSaveDays)
//	Logger.Infof("loginit success ! \n")
//
//
//}

func main() {


	log2.Logger.Info("main startlll: \n")
	log2.Logger.Debug("hello worldlll! \n")
	log2.Logger.Warn("log2 warning ")
	log2.Logger.Warn("log2 warninglll ")
	log2.Logger.Warn("log2 warninglll ")
	log2.Logger.Warn("log2 warninglll ")
	log2.Logger.Warn("log2 warninglll ")
	log2.Logger.Error("hello worldlll!")
	log2.Logger.Fatal("hello worldlll!")

	//logs.Errorf("-- %s\n", "hello world!")

	//log3.Debug("log3 debug")
	//log3.Info("log3 info")
	//log3.Error("log3 error")
	//log3.Warn("log3 warnning....")
}
