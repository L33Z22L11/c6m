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

	db.StoreMsg(msg)

	// 在线列表
	// if model.Connections[msg.Dest] == nil {
	// 	db.AddMsg()
	// 	return
	// }

	var dest []string
	if msg.Src == "groupReq" {
		adminMap, _ := db.ListGroupAdmin(msg.Dest)
		for k := range adminMap {
			dest = append(dest, k)
		}
	} else if msg.Dest[0] == 'g' {
		destMap, _ := db.ListGroupMember(msg.Dest)
		for k := range destMap {
			dest = append(dest, k)
		}
	} else {
		dest = append(dest, msg.Dest)
	}

	for _, d := range dest {
		if conn, ok := model.Connections[d]; ok {
			conn.WriteJSON(msg)
		}
	}
}

func IsVaildMsg(msg *model.Message) bool {
	if msg.Src[0] == 'u' {
		if msg.Dest[0] == 'u' {
			return db.IsFriend(msg.Src, msg.Dest)
		} else if msg.Dest[0] == 'g' {
			return db.IsGroupMember(msg.Src, msg.Dest)
		}
	}
	return true
}

// 登录时检查离线消息队列
func CheckMsg(msg *model.Message) {

}
