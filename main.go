package main

import (
	"c6m/db"
	"c6m/router"
	"os"
	"os/signal"
)

func main() {
	// 创建一个信号通道
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	db.StartServer()
	db.Test()

	go router.InitWebServer()

	// 等待信号
	<-sigCh
	db.StopServer()
	db.Test()
}
