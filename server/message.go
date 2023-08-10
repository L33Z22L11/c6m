package server

import (
	"c6m/db"
	"c6m/model"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// 在线推送
func PushMsg(msg *model.Message) {
	// 打印调试信息
	msgJson, _ := json.Marshal(msg)
	fmt.Printf("%s\n", msgJson)

	if !IsVaildMsg(msg) {
		msg = &model.Message{
			Type:    "toast",
			Src:     "0",
			Dest:    msg.Src,
			Content: fmt.Sprintf("没有权限发送这条消息：%s", msgJson),
		}
	}
	if msg.Dest[0] == '-' {
		panic("暂不支持群聊推送")
	}

	msg.Time = strconv.FormatInt(time.Now().UnixMicro(), 10)
	if model.Connections[msg.Dest] == nil {
		db.AddMsg()
		return
	}

	model.Connections[msg.Dest].WriteJSON(msg)
}

func IsVaildMsg(msg *model.Message) bool {
	if msg.Src == "0" {
		return true
	}
	switch msg.Type {
	case "single":
		return db.IsFriend(msg.Src, msg.Dest)
	case "group":
	}
	return false
}

// 登录时检查离线消息队列
func CheckMsg(msg *model.Message) {

}
