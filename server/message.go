package server

import (
	"c6m/db"
	"c6m/model"
	"encoding/json"
	"fmt"
	"time"
)

// 在线推送
func PushMsg(msg *model.Message) {
	// 打印调试信息
	msgJson, _ := json.Marshal(msg)
	fmt.Printf("%s\n", msgJson)

	if !IsVaildMsg(msg) {
		msg = &model.Message{
			Src:     "system",
			Dest:    msg.Src,
			Content: fmt.Sprintf("没有权限发送这条消息：%s", msgJson),
		}
	}

	msg.Time = time.Now().UnixNano() / int64(time.Millisecond)

	if model.Connections[msg.Dest] == nil {
		db.AddMsg()
		return
	}

	model.Connections[msg.Dest].WriteJSON(msg)
}

func IsVaildMsg(msg *model.Message) bool {
	if msg.Src[0] == 'u' {
		if msg.Dest[0] == 'u' {
			return db.HaveFriend(msg.Src, msg.Dest)
		} else if msg.Dest[0] == 'g' {
			return db.HaveGroup(msg.Src, msg.Dest)
		}
	}
	return true
}

// 登录时检查离线消息队列
func CheckMsg(msg *model.Message) {

}
