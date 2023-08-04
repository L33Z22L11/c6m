package db

import (
	"c6m/model"
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
	server.PushMsg(&model.Message{})

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

	isFriend, _ := rc.SIsMember(context.Background(), fmt.Sprintf("friend:%s", friendUid), uid).Result()
	if !isFriend {
		return fmt.Errorf("还不是对方好友")
	}

	rc.HDel(context.Background(), fmt.Sprintf("friend:%s", uid), friendUid)
	rc.HDel(context.Background(), fmt.Sprintf("friend:%s", friendUid), uid)

	return nil
}

func GetFriendReq(uid string) ([]string, error) {
	friendReqList, err := rc.SMembers(context.Background(), fmt.Sprintf("friendReq:%s", uid)).Result()
	if err != nil {
		return nil, fmt.Errorf("获取好友请求列表失败:%s", err)
	}

	return friendReqList, err
}

func RespFriendReq(uid string, friendUid string, isAccept string) error {
	havefriendReq, _ := rc.SIsMember(context.Background(), fmt.Sprintf("friendReq:%s", uid), friendUid).Result()
	if !havefriendReq {
		return fmt.Errorf("不存在这个好友申请")
	}

	if isAccept == "1" {
		rc.SAdd(context.Background(), fmt.Sprintf("friend:%s", uid), friendUid)
		rc.SAdd(context.Background(), fmt.Sprintf("friend:%s", friendUid), uid)
	}

	rc.SRem(context.Background(), fmt.Sprintf("friendReq:%s", uid), friendUid)
	return nil
}

func ListFriend(uid string) ([]string, error) {
	friendList, err := rc.SMembers(context.Background(), fmt.Sprintf("friend:%s", uid)).Result()
	if err != nil {
		return nil, fmt.Errorf("获取好友列表失败:%s", err)
	}

	return friendList, err
}
