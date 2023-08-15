package router

import (
	"c6m/db"
	"c6m/model"
	"c6m/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleWebSocket(c *gin.Context) {
	uid, err := db.GetUidByToken(c.Query("token"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		fmt.Printf("token:%s uid:%s err:%s\n", c.Query("token"), uid, err.Error())
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket升级失败: ", err)
		return
	}

	// 阻止重复登录
	if model.Connections[uid] != nil {
		model.Connections[uid].Close()
		delete(model.Connections, uid)
	}
	model.Connections[uid] = conn

	msg := model.Message{
		Src:     "system",
		Dest:    uid,
		Content: "欢迎用户" + db.MustGetNameById(uid),
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
		msg.Src = uid
		server.PushMsg(&msg)
	}

	// 关闭WebSocket连接
	conn.Close()
	delete(model.Connections, uid)
}
