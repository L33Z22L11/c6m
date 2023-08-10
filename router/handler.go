package router

import (
	"c6m/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var router = gin.Default()
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域连接
	},
}

func InitWebServer() {
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)

	model.Connections = make(map[string]*websocket.Conn)
	router.GET("/ws", handleWebSocket)

	router.POST("/friend/add", VerifyToken(), handleAddFriend)
	router.POST("/friend/del", VerifyToken(), handleDelFriend)
	router.GET("/friend/req", VerifyToken(), handleGetFriendReq)
	router.POST("/friend/req", VerifyToken(), handleRespFriendReq)
	router.GET("/friend/all", VerifyToken(), handleListFriend)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/app")
	})
	router.Static("/app", "./web")

	err := router.Run(":4000")
	if err != nil {
		fmt.Printf("网页服务器启动失败: %s\n", err)
		return
	}
	fmt.Println("网页服务器启动")
}
