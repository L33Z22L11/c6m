package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func AddFriend(uid string, friendName string, content string) (string, error) {
	friendUid, err := GetUidByUname(friendName)
	if err != nil {
		return "", err
	}

	if uid == friendUid {
		return "", fmt.Errorf("不能添加自己为好友")
	}

	_, err = rc.HGet(context.Background(), "friend:"+uid, friendUid).Result()
	if err != redis.Nil {
		return "", fmt.Errorf("已经是好友了")
	}

	requested, _ := rc.HGet(context.Background(), "friendReq:"+friendUid, uid).Result()
	if requested != "" {
		return "", fmt.Errorf("已经发送过一次好友申请")
	}

	friendReqList, _ := GetFriendReq(uid)
	if friendReqList[friendUid] != "" {
		RespFriendReq(uid, friendUid, "1")
	} else {
		rc.HSet(context.Background(), "friendReq:"+friendUid, uid, content)
	}
	return friendUid, nil
}

func DelFriend(uid string, friendName string) error {
	friendUid, err := GetUidByUname(friendName)
	if err != nil {
		return err
	}

	if uid == friendUid {
		return fmt.Errorf("不能删除自己")
	}

	haveFriend, _ := rc.SIsMember(context.Background(), "friend:"+friendUid, uid).Result()
	if !haveFriend {
		return fmt.Errorf("还不是对方好友")
	}

	rc.SRem(context.Background(), "friend:"+uid, friendUid)
	rc.SRem(context.Background(), "friend:"+friendUid, uid)

	return nil
}

func GetFriendReq(uid string) (map[string]string, error) {
	friendReqList, err := rc.HGetAll(context.Background(), "friendReq:"+uid).Result()
	if err != nil {
		return nil, fmt.Errorf("获取好友请求列表失败:%s", err)
	}

	return friendReqList, err
}

func RespFriendReq(uid string, friendUid string, isAccept string) error {
	_, err := rc.HGet(context.Background(), "friendReq:"+uid, friendUid).Result()
	if err == redis.Nil {
		return fmt.Errorf("不存在这个好友申请")
	}

	if isAccept == "1" {
		rc.SAdd(context.Background(), "friend:"+uid, friendUid)
		rc.SAdd(context.Background(), "friend:"+friendUid, uid)
	}

	rc.HDel(context.Background(), "friendReq:"+uid, friendUid)
	return nil
}

func ListFriend(uid string) (map[string]string, error) {
	friendList, err := rc.SMembers(context.Background(), "friend:"+uid).Result()
	if err != nil {
		return nil, fmt.Errorf("获取好友列表失败:%s", err)
	}

	friendMap := make(map[string]string)

	for _, friendUid := range friendList {
		friendMap[friendUid] = MustGetNameById(friendUid)
	}
	return friendMap, err
}

func HaveFriend(uid string, friendUid string) bool {
	// 查询发送者的好友列表
	friendMap, _ := ListFriend(uid)

	// 判断接收者是否在好友列表中
	return friendMap[friendUid] != ""
}
