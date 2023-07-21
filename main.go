package main

import (
	db "c6m/database"
	"c6m/util"
	"os"
	"os/signal"
)

func main() {
	// 创建一个信号通道
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	db.StartServer()
	db.Test()

	go util.InitWebServer()

	// ...

	// 等待信号
	<-sigCh
	db.StopServer()
	db.Test()
}
