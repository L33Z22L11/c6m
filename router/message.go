package router

import (
	"c6m/db"
	"c6m/model"
	"c6m/server"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handleGetHistory(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	id := c.Query("id")

	msgStorage, err := db.GetMsgStorage(uid, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, msgStorage)
}

func handleUpload(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	dest := c.PostForm("dest")
	fileName := fmt.Sprintf("%s_%s_%s", uid, time.Now().Format("20060102_150405"), file.Filename)

	// 将文件保存到公共目录
	err := c.SaveUploadedFile(file, "public/"+fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "上传文件失败",
		})
		return
	}

	// 设置文件过期时间为3天后
	expirationTime := time.Now().Add(3 * 24 * time.Hour)

	msg := model.Message{
		Src:     uid,
		Dest:    dest,
		Content: fmt.Sprintf("用户%s发送了<a href='/public/%s'>%s</a>(于%s过期)", db.MustGetNameById(uid), fileName, file.Filename, expirationTime.Format("2006-01-02 15:04")),
	}

	c.JSON(http.StatusOK, gin.H{
		"src":     uid,
		"dest":    dest,
		"time":    time.Now().UnixNano() / int64(time.Millisecond),
		"content": fmt.Sprintf("用户%s发送了<a href='/public/%s'>%s</a>(于%s过期)", db.MustGetNameById(uid), fileName, file.Filename, expirationTime.Format("2006-01-02 15:04")),
	})
	server.PushMsg(&msg)
}
