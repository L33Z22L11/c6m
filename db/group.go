package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func AddGroup(uid string, groupName string, content string) (string, error) {
	gid, err := GetGidByGname(groupName)
	if err != nil {
		return "", err
	}

	requested, _ := rc.HGet(context.Background(), fmt.Sprintf("groupReq:%s", gid), uid).Result()
	if requested != "" {
		return "", fmt.Errorf("已经发送过一次群申请")
	}

	groupReqList, _ := GetGroupReq(uid)
	if groupReqList[gid] != "" {
		RespGroupReq(uid, gid, "1")
	} else {
		rc.HSet(context.Background(), fmt.Sprintf("groupReq:%s", gid), uid, content)
	}
	return gid, nil
}

func DelGroup(gid string, groupName string) error {
	groupgid, err := GetGidByGname(groupName)
	if err != nil {
		return err
	}

	if gid == groupgid {
		return fmt.Errorf("不能删除自己")
	}

	isGroup, _ := rc.SIsMember(context.Background(), fmt.Sprintf("group:%s", groupgid), gid).Result()
	if !isGroup {
		return fmt.Errorf("还不是对方群")
	}

	rc.SRem(context.Background(), fmt.Sprintf("group:%s", gid), groupgid)
	rc.SRem(context.Background(), fmt.Sprintf("group:%s", groupgid), gid)

	return nil
}

func GetGroupReq(gid string) (map[string]string, error) {
	groupReqList, err := rc.HGetAll(context.Background(), fmt.Sprintf("groupReq:%s", gid)).Result()
	if err != nil {
		return nil, fmt.Errorf("获取群请求列表失败:%s", err)
	}

	return groupReqList, err
}

func RespGroupReq(gid string, groupgid string, isAccept string) error {
	_, err := rc.HGet(context.Background(), fmt.Sprintf("groupReq:%s", gid), groupgid).Result()
	if err == redis.Nil {
		return fmt.Errorf("不存在这个群申请")
	}

	if isAccept == "1" {
		rc.SAdd(context.Background(), fmt.Sprintf("group:%s", gid), groupgid)
		rc.SAdd(context.Background(), fmt.Sprintf("group:%s", groupgid), gid)
	}

	rc.HDel(context.Background(), fmt.Sprintf("groupReq:%s", gid), groupgid)
	return nil
}

func ListGroup(gid string) (map[string]string, error) {
	groupList, err := rc.SMembers(context.Background(), fmt.Sprintf("group:%s", gid)).Result()
	if err != nil {
		return nil, fmt.Errorf("获取群列表失败:%s", err)
	}

	groupMap := make(map[string]string)

	for _, groupgid := range groupList {
		groupMap[groupgid] = MustGetNameById(groupgid)
	}
	return groupMap, err
}

func HaveGroup(gid string, groupgid string) bool {
	// 查询发送者的群列表
	groupMap, _ := ListGroup(gid)

	// 判断接收者是否在群列表中
	return groupMap[groupgid] != ""
}
