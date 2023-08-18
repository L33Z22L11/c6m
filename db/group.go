package db

import (
	"context"
	"fmt"
	"strings"
)

func JoinGroup(uid string, gname string, content string) (string, error) {
	gid, err := GetGidByGname(gname)
	if err != nil {
		return "", err
	}

	if IsGroupMember(uid, gid) {
		return "", fmt.Errorf("已经是群成员")
	}

	requested, _ := rc.HGet(context.Background(), "groupReq:"+gid, uid).Result()
	if requested != "" {
		return "", fmt.Errorf("已经发送过一次群申请")
	}

	rc.HSet(context.Background(), "groupReq:"+gid, uid, content)
	return gid, nil
}

func LeaveGroup(uid string, gname string) error {
	gid, err := GetGidByGname(gname)
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

	rc.SRem(context.Background(), "groupAdmin:"+gid, uid)
	rc.SRem(context.Background(), "group:"+gid, uid)
	rc.SAdd(context.Background(), "inGroup:"+uid, gid)

	return nil
}

func InfoGroup(uid, gname string) (string, error) {
	gid, err := GetGidByGname(gname)
	if err != nil {
		return "", fmt.Errorf("不存在这个群")
	}
	if !IsGroupMember(uid, gid) {
		return "", fmt.Errorf("您不是该群的成员")
	}

	group, _ := GetGroupByGid(gid)
	adminList, _ := rc.SMembers(context.Background(), "groupAdmin:"+gid).Result()
	memberList, _ := rc.SMembers(context.Background(), "group:"+gid).Result()

	info := fmt.Sprintf("gid: %s, 群主: %s, 管理员: %s, 成员: %s",
		gid, group.Owner, strings.Join(adminList, " & "), strings.Join(memberList, " & "))

	return info, nil
}

func GetGroup(uid string) (map[string]string, error) {
	groupList, err := rc.SMembers(context.Background(), "inGroup:"+uid).Result()
	if err != nil {
		return nil, fmt.Errorf("获取群成员列表失败:%s", err)
	}

	groupMap := make(map[string]string)

	for _, gid := range groupList {
		groupMap[gid] = MustGetNameById(gid)
	}
	return groupMap, err
}

func ListGroupMember(gid string) (map[string]string, error) {
	memberList, err := rc.SMembers(context.Background(), "group:"+gid).Result()
	if err != nil {
		return nil, fmt.Errorf("获取群成员列表失败:%s", err)
	}

	groupMap := make(map[string]string)

	for _, gid := range memberList {
		groupMap[gid] = MustGetNameById(gid)
	}
	return groupMap, err
}

func IsGroupMember(uid string, gid string) bool {
	// 查询发送者的群列表
	groupMap, _ := ListGroupMember(gid)

	// 判断接收者是否在群列表中
	return groupMap[uid] != ""
}
