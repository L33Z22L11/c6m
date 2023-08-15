package router

import (
	"c6m/db"
	"c6m/model"
	"c6m/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleAddGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")
	content := c.PostForm("content")
	if content == "" {
		content = "我是" + db.MustGetNameById(uid)
	}

	groupGuid, err := db.JoinGroup(uid, groupName, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	server.PushMsg(&model.Message{
		Src:     "groupReq",
		Dest:    groupGuid,
		Content: db.MustGetNameById(uid) + "请求加群" + groupName,
	})
	c.JSON(http.StatusOK, gin.H{
		"message":    "已发送好友申请",
		"group_name": groupName,
	})
}

func handleDelGroup(c *gin.Context) {
	guid := c.MustGet("guid").(string)
	groupName := c.PostForm("group_name")

	err := db.LeaveGroup(guid, groupName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "已删除好友",
		"group_name": groupName,
	})
}

func handleGetGroupReq(c *gin.Context) {
	guid := c.MustGet("guid").(string)

	groupReqList, err := db.GetGroupReq(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groupReqList)
}

func handleRespGroupReq(c *gin.Context) {
	guid := c.MustGet("guid").(string)
	groupGuid := c.PostForm("group_guid")
	isAccept := c.PostForm("accept")

	err := db.RespGroupReq(guid, groupGuid, isAccept)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var content string
	if isAccept == "1" {
		content = "成功加入群" + db.MustGetNameById(guid)
	} else {
		content = db.MustGetNameById(guid) + "拒绝让你加入"
	}

	server.PushMsg(&model.Message{
		Src:     "groupReq",
		Dest:    groupGuid,
		Content: content,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":    "已处理请求",
		"group_name": db.MustGetNameById(guid),
	})
}

func handleListGroup(c *gin.Context) {
	guid := c.MustGet("guid").(string)

	groupList, err := db.ListGroupMembers(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groupList)
}
