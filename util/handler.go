package util

import (
	db "c6m/database"
	"fmt"
	"log"
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
	router.GET("/ws", VerifyToken(), handleWebSocket)
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)

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

// 处理注册请求
func handleRegister(c *gin.Context) {
	// 获取请求参数
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 用户注册逻辑
	user, err := db.CreateUser(username, password)

	// 返回响应
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"uid":      user.Uid,
		"username": user.Username,
	})
}

// 处理登录请求
func handleLogin(c *gin.Context) {
	// 获取请求参数
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 用户登录逻辑
	token, err := db.AuthUser(username, password)

	// 返回响应
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"token":    token,
	})
}

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		uid, err := db.GetUidByToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Set("uid", uid)
		c.Next()
	}
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket升级失败: ", err)
		return
	}

	// 在这里可以处理WebSocket连接
	uid := c.MustGet("uid").(string)
	auth, _ := db.GetAuthByUID(uid)
	conn.WriteJSON(&Message{
		Type: "toast",
		Text: fmt.Sprintf("欢迎用户%s", auth.Username),
	})

	// 读取和处理来自客户端的消息
	for {
		// 读取消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket读消息失败: ", err)
			break
		}

		// 在这里可以处理接收到的消息
		// ...

		// 发送响应消息
		err = conn.WriteMessage(websocket.TextMessage, []byte("收到消息: "+string(msg)))
		if err != nil {
			log.Println("WebSocket响应失败: ", err)
			break
		}
	}

	// 关闭WebSocket连接
	conn.Close()
}

func handleAddFriend(c *gin.Context) {
	uid := c.MustGet("uid").(string)

	friendName := c.PostForm("friend_name")

	err := db.AddFriend(uid, friendName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "已发送好友申请",
		"uid":         uid,
		"friend_name": friendName,
	})
}
