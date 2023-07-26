package database

import (
	"c6m/models"
	"c6m/server"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func AddFriend(uid string, friendName string) error {
	friendUid, err := GetUidByUname(friendName)
	if err != nil {
		return err
	}

	if uid == friendUid {
		return fmt.Errorf("不能添加自己为好友")
	}

	_, err = rc.HGet(context.Background(), fmt.Sprintf("friend:%s", uid), friendUid).Result()
	if err != redis.Nil {
		return fmt.Errorf("已经是好友了")
	}

	requested, _ := rc.SIsMember(context.Background(), fmt.Sprintf("friendReq:%s", friendUid), uid).Result()
	if requested {
		return fmt.Errorf("已经发送过一次好友申请")
	}

	rc.SAdd(context.Background(), fmt.Sprintf("friendReq:%s", friendUid), uid)
	server.ParseMsg(&models.Message{})

	return nil
}

func DelFriend(uid string, friendName string) error {
	friendUid, err := GetUidByUname(friendName)
	if err != nil {
		return err
	}

	if uid == friendUid {
		return fmt.Errorf("不能删除自己")
	}

	return nil
}
