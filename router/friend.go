package router

import (
	"c6m/db"
	"c6m/model"
	"c6m/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAddFriend(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	friendName := c.PostForm("friend_name")
	content := c.PostForm("content")
	if content == "" {
		content = "我是" + db.MustGetNameById(uid)
	}

	friendUid, err := db.AddFriend(uid, friendName, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	server.PushMsg(&model.Message{
		Src:     "friendReq",
		Dest:    friendUid,
		Content: db.MustGetNameById(uid) + "请求添加你为好友",
	})
	c.JSON(http.StatusOK, gin.H{
		"message":     "已发送好友申请",
		"friend_name": friendName,
	})
}

func handleDelFriend(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	friendName := c.PostForm("friend_name")

	err := db.DelFriend(uid, friendName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "已删除好友",
		"friend_name": friendName,
	})
}

func handleGetFriendReq(c *gin.Context) {
	uid := c.MustGet("uid").(string)

	friendReqList, err := db.GetFriendReq(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, friendReqList)
}

func handleRespFriendReq(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	friendUid := c.PostForm("friend_uid")
	isAccept := c.PostForm("accept")

	err := db.RespFriendReq(uid, friendUid, isAccept)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var content string
	if isAccept == "1" {
		content = db.MustGetNameById(uid) + "通过了你的好友申请"
	} else {
		content = db.MustGetNameById(uid) + "拒绝了你的好友申请"
	}

	server.PushMsg(&model.Message{
		Src:     "friendReq",
		Dest:    friendUid,
		Content: content,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":     "已处理请求",
		"friend_name": db.MustGetNameById(uid),
	})
}

func handleListFriend(c *gin.Context) {
	uid := c.MustGet("uid").(string)

	friendList, err := db.ListFriend(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, friendList)
}
