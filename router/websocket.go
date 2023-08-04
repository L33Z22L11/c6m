package router

import (
	"c6m/db"
	"c6m/model"
	"c6m/server"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket升级失败: ", err)
		return
	}

	// 在这里可以处理WebSocket连接
	uid := c.MustGet("uid").(string)
	msg := model.Message{
		Type:    "toast",
		Content: fmt.Sprintf("欢迎用户%s", db.MustGetUnameByUID(uid)),
	}
	server.PushMsg(&msg)

	// 读取和处理来自客户端的消息
	for {
		// 读取消息
		_, msgJson, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket读消息失败: ", err)
			break
		}

		json.Unmarshal(msgJson, &msg)
		server.PushMsg(&msg)
	}

	// 关闭WebSocket连接
	conn.Close()

}
