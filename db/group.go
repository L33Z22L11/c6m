package db

import (
	"context"
	"fmt"
)

func JoinGroup(uid string, groupName string, content string) (string, error) {
	gid, err := GetGidByGname(groupName)
	if err != nil {
		return "", err
	}

	requested, _ := rc.HGet(context.Background(), "groupReq:"+gid, uid).Result()
	if requested != "" {
		return "", fmt.Errorf("已经发送过一次群申请")
	}

	groupReqList, _ := GetGroupReq(uid)
	if groupReqList[gid] != "" {
		RespGroupReq(uid, gid, "1")
	} else {
		rc.HSet(context.Background(), "groupReq:"+gid, uid, content)
	}
	return gid, nil
}

func LeaveGroup(uid string, groupName string) error {
	gid, err := GetGidByGname(groupName)
	if err != nil {
		return err
	}

	haveGroup, _ := rc.SIsMember(context.Background(), "group:"+gid, uid).Result()
	if !haveGroup {
		return fmt.Errorf("未加入此群")
	}

	group, _ := GetGroupByGid(gid)
	if group.Owner == uid {
		return fmt.Errorf("删除失败，群主请使用解散群功能")
	}

	rc.SRem(context.Background(), "gadmin:"+gid, uid)
	rc.SRem(context.Background(), "group:"+gid, uid)

	return nil
}

func ListGroupMembers(gid string) (map[string]string, error) {
	memberList, err := rc.SMembers(context.Background(), "group:"+gid).Result()
	if err != nil {
		return nil, fmt.Errorf("获取群成员列表失败:%s", err)
	}

	groupMap := make(map[string]string)

	for _, groupgid := range memberList {
		groupMap[groupgid] = MustGetNameById(groupgid)
	}
	return groupMap, err
}

func HaveGroup(gid string, groupgid string) bool {
	// 查询发送者的群列表
	groupMap, _ := ListGroupMembers(gid)

	// 判断接收者是否在群列表中
	return groupMap[groupgid] != ""
}
