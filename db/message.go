package db

import (
	"c6m/model"
	"context"
	"encoding/json"
)

func AddMsg(msg *model.Message) {
	// 序列化Message结构体为JSON
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		// 处理序列化错误
		return
	}

	// 将JSON写入Redis队列
	err = rc.RPush(context.Background(), "message:"+msg.Dest, jsonMsg).Err()
	if err != nil {
		// 处理写入Redis错误
		return
	}
}
