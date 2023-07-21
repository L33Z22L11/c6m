package util

import (
	db "c6m/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func InitWebServer() {
	router.Static("/", "./web")
	router.POST("/register", handleRegister)
	router.POST("/login", handleLogin)
	router.POST("/add-friend", handlerAddFriend)

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
		c.String(http.StatusBadRequest, err.Error())
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
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"token":    token,
	})
}

func handlerAddFriend(c *gin.Context) {
	username := c.PostForm("username")
	friendName := c.PostForm("friend_name")

	err := db.AddFriend(username, friendName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":   username,
		"friendName": friendName,
	})
}
