package router

import (
	"c6m/db"
	"c6m/model"
	"c6m/server"
	"fmt"
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

	groupgid, err := db.JoinGroup(uid, groupName, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	server.PushMsg(&model.Message{
		Src:     "groupReq",
		Dest:    groupgid,
		Content: db.MustGetNameById(uid) + "请求加群" + groupName,
	})
	c.JSON(http.StatusOK, gin.H{
		"message":    "已发送加群申请",
		"group_name": groupName,
	})
}

func handleLeaveGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")

	err := db.LeaveGroup(uid, groupName)
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

	err := db.InviteGroup(uid, groupName, invitee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

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

	err := db.KickGroup(uid, groupName, member)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "已将成员踢出群组",
		"group_name": groupName,
		"member":     member,
	})
}

func handleCreateGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")

	group, err := db.CreateGroup(uid, groupName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "群组创建成功",
		"gid":        group.Gid,
		"group_name": group.Gname,
		"owner":      group.Owner,
	})
}

func handleDelGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")

	err := db.DelGroup(uid, groupName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "群组已解散",
		"group_name": groupName,
	})
}

func handleAddGadmin(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	groupName := c.PostForm("group_name")
	admin := c.PostForm("admin")

	err := db.AddGroupAdmin(uid, groupName, admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

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

	err := db.DelGroupAdmin(uid, groupName, admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "已移除群组管理员",
		"group_name": groupName,
		"admin":      admin,
	})
}

func handleGetGroupReq(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	gid := c.Query("gid")

	groupReqList, err := db.GetGroupReq(uid, gid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groupReqList)
}

func handleRespGroupReq(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	gid := c.PostForm("gid")
	isAccept := c.PostForm("accept")

	err := db.RespGroupReq(uid, gid, isAccept)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var content string
	if isAccept == "1" {
		content = "成功加入群" + db.MustGetNameById(gid)
	} else {
		content = db.MustGetNameById(gid) + "拒绝让你加入"
	}

	server.PushMsg(&model.Message{
		Src:     "groupReq",
		Dest:    gid,
		Content: content,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":    "已处理请求",
		"group_name": db.MustGetNameById(gid),
	})
}

func handleGetGroup(c *gin.Context) {
	uid := c.MustGet("uid").(string)

	groupList, err := db.GetGroup(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groupList)
}

func handleListGroupMember(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	gid := c.Query("gid")

	if !db.IsGroupMember(uid, gid) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("您不是该群的成员"),
		})
		return
	}

	groupList, err := db.ListGroupMember(gid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, groupList)
}
