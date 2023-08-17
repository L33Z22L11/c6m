package router

import (
	"c6m/db"
	"c6m/model"
	"c6m/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleJoinGroup(c *gin.Context) {
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
		"message":    "已发送加群申请",
		"group_name": groupName,
	})
}

func handleLeaveGroup(c *gin.Context) {
	guid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")

	err := db.LeaveGroup(guid, groupName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "已退群",
		"group_name": groupName,
	})
}

func handleInviteGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")
	invitee := c.PostForm("invitee")

	// TODO: Implement logic for handling group invitation

	c.JSON(http.StatusOK, gin.H{
		"message":    "邀请已发送",
		"group_name": groupName,
		"invitee":    invitee,
	})
}

func handleKickGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")
	member := c.PostForm("member")

	// TODO: Implement logic for handling kicking a member from a group

	c.JSON(http.StatusOK, gin.H{
		"message":    "已将成员踢出群组",
		"group_name": groupName,
		"member":     member,
	})
}

func handleCreateGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")

	// TODO: Implement logic for handling group creation

	c.JSON(http.StatusOK, gin.H{
		"message":    "群组创建成功",
		"group_name": groupName,
	})
}

func handleDelGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")

	// TODO: Implement logic for handling group dissolution

	c.JSON(http.StatusOK, gin.H{
		"message":    "群组已解散",
		"group_name": groupName,
	})
}

func handleAddGadmin(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")
	admin := c.PostForm("admin")

	// TODO: Implement logic for handling adding a group admin

	c.JSON(http.StatusOK, gin.H{
		"message":    "已添加群组管理员",
		"group_name": groupName,
		"admin":      admin,
	})
}

func handleDelGadmin(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")
	admin := c.PostForm("admin")

	// TODO: Implement logic for handling removing a group admin

	c.JSON(http.StatusOK, gin.H{
		"message":    "已移除群组管理员",
		"group_name": groupName,
		"admin":      admin,
	})
}

func handleGetGroupReq(c *gin.Context) {
	guid := c.MustGet("uid").(string)

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
	guid := c.MustGet("uid").(string)
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
	guid := c.MustGet("uid").(string)

	groupList, err := db.ListGroupMembers(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groupList)
}
