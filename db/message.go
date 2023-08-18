package db

import (
	"c6m/model"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func StoreMsg(msg *model.Message) error {
	// 序列化Message结构体为JSON
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		// 处理序列化错误
		return err
	}

	msgStorage := "message:"

	switch msg.Dest[0] {
	case 'u':
		uids := []string{msg.Src, msg.Dest}
		sort.Strings(uids)
		msgStorage += strings.Join(uids, ":")
	case 'g':
		msgStorage += msg.Dest
	default:
		return fmt.Errorf("无效的msg.Dest:%s", msg.Dest)
	}

	// 将JSON写入Redis队列
	err = rc.RPush(context.Background(), msgStorage, jsonMsg).Err()
	if err != nil {
		// 处理写入Redis错误
		return err
	}
	return nil
}

func GetMsgStorage(uid, id string) ([]*model.Message, error) {
	var storageKey string
	switch id[0] {
	case 'u':
		if !IsFriend(uid, id) {
			return nil, fmt.Errorf("不是好友:%s", id)
		}
		uids := []string{uid, id}
		sort.Strings(uids)
		storageKey = "message:" + strings.Join(uids, ":")
	case 'g':
		if !IsGroupMember(uid, id) {
			return nil, fmt.Errorf("不在此群:%s", id)
		}
		storageKey = "message:" + id

	default:
		return nil, fmt.Errorf("未知id:%s", id)
	}

	// 从Redis队列中读取JSON消息
	jsonMsgs, err := rc.LRange(context.Background(), storageKey, 0, -1).Result()
	if err != nil {
		// 处理从Redis读取错误
		return nil, err
	}

	// 反序列化JSON消息为Message结构体
	var messages []*model.Message
	for _, jsonMsg := range jsonMsgs {
		var msg model.Message
		err := json.Unmarshal([]byte(jsonMsg), &msg)
		if err != nil {
			// 处理反序列化错误
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}
