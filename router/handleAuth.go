package router

import (
	"c6m/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 处理注册请求
func handleRegister(c *gin.Context) {
	// 获取请求参数
	username := c.PostForm("username")
	password := c.PostForm("password")
	question := c.PostForm("question")
	answer := c.PostForm("answer")

	// 用户注册逻辑
	user, err := db.CreateUser(username, password, question, answer)

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
	uid, token, err := db.AuthUser(username, password)

	// 返回响应
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, gin.H{
		"uid":      uid,
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
			c.Abort()
			return
		}
		c.Set("uid", uid)
		c.Next()
	}
}

func handleChangePw(c *gin.Context) {
	// 获取请求参数
	username := c.PostForm("username")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")

	// 调用 ChangePassword 函数修改密码
	err := db.ChangePassword(username, oldPassword, newPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

func handleAuthQuestion(c *gin.Context) {
	// 获取请求参数
	username := c.Query("username")

	// 调用 GetAuthQuestionByUname 函数获取安全问题
	question, err := db.GetAuthQuestionByUname(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"question": question})
}

func handleResetPw(c *gin.Context) {
	// 获取请求参数
	username := c.PostForm("username")
	answer := c.PostForm("answer")
	newPassword := c.PostForm("new_password")

	// 调用 ResetPwByAnswer 函数重置密码
	err := db.ResetPwByAnswer(username, answer, newPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
