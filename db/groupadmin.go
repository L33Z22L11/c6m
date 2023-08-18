package db

import (
	"c6m/model"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func CreateGroup(uid, gname string) (*model.Group, error) {
	if len(gname) == 0 {
		return nil, fmt.Errorf("群名称不能为空")
	}

	gid, _ := GetGidByGname(gname)
	if gid != "" {
		return nil, fmt.Errorf("群已存在")
	}
	group := &model.Group{
		Gid:   generateGid(),
		Gname: gname,
		Owner: uid,
	}

	err := SaveGroup(group)
	if err != nil {
		return nil, fmt.Errorf("保存群[%s]到Redis: %v", gname, err)
	}

	rc.SAdd(context.Background(), "group:"+group.Gid, uid)
	rc.SAdd(context.Background(), "inGroup:"+uid, group.Gid)

	return group, nil
}

func SaveGroup(group *model.Group) error {
	// 将群组信息序列化为 JSON 字符串
	groupJson, err := json.Marshal(group)
	if err != nil {
		return err
	}

	// 将群组信息存储到 Redis 哈希表中
	err = rc.HSet(context.Background(), "group", group.Gid, groupJson).Err()
	if err != nil {
		return err
	}
	err = rc.HSet(context.Background(), "gid", group.Gname, group.Gid).Err()
	if err != nil {
		return err
	}

	return nil
}

func DelGroup(uid, gname string) error {
	// Check if the group exists
	gid, err := GetGidByGname(gname)
	if err != nil {
		return fmt.Errorf("群组[%s]不存在", gid)
	}

	if !IsGroupOwner(uid, gid) {
		return fmt.Errorf("不是群主")
	}

	// Delete the group from Redis
	kickAll(gid)

	err = rc.HDel(context.Background(), "group", gid).Err()
	if err != nil {
		return fmt.Errorf("从Redis中删除群组[%s]出错: %v", gid, err)
	}
	err = rc.HDel(context.Background(), "gid", gname, gid).Err()
	if err != nil {
		return fmt.Errorf("从Redis中删除群组[%s]出错: %v", gid, err)
	}

	return nil
}

func AddGroupAdmin(uid, gname, admin string) error {
	gid, err := GetGidByGname(gname)
	if err != nil {
		return fmt.Errorf("获取群组信息失败：%v", err)
	}

	if !IsGroupMember(uid, gid) {
		return fmt.Errorf("您不是该群的成员")
	}

	if !IsGroupOwner(uid, gid) {
		return fmt.Errorf("您不是群主")
	}

	if IsGroupAdmin(admin, gid) {
		return fmt.Errorf("该用户已是群管理员")
	}

	err = rc.SAdd(context.Background(), "groupAdmin:"+admin, gid).Err()
	if err != nil {
		return fmt.Errorf("添加群管理员失败：%v", err)
	}

	return nil
}

func DelGroupAdmin(uid, gname, admin string) error {
	gid, err := GetGidByGname(gname)
	if err != nil {
		return fmt.Errorf("获取群组信息失败：%v", err)
	}

	if !IsGroupMember(uid, gid) {
		return fmt.Errorf("您不是该群的成员")
	}

	if !IsGroupOwner(uid, gid) {
		return fmt.Errorf("您不是群主")
	}

	if !IsGroupAdmin(admin, gid) {
		return fmt.Errorf("该用户不是群管理员")
	}

	err = rc.SRem(context.Background(), "groupAdmin:"+admin, gid).Err()
	if err != nil {
		return fmt.Errorf("移除群管理员失败：%v", err)
	}

	return nil
}

func KickGroup(uid, gname, member string) error {
	gid, err := GetGidByGname(gname)
	if err != nil {
		return fmt.Errorf("获取群组信息失败：%v", err)
	}

	if !IsGroupMember(uid, gid) {
		return fmt.Errorf("您不是该群的成员")
	}

	if !IsGroupMember(gid, member) {
		return fmt.Errorf("该用户不是群成员")
	}

	if IsGroupAdmin(gid, member) && !IsGroupAdmin(uid, gid) && !IsGroupOwner(uid, gid) {
		return fmt.Errorf("您没有权限踢出群管理员/群主")
	}

	err = rc.SRem(context.Background(), "group:"+gid, member).Err()
	if err != nil {
		return fmt.Errorf("踢出群成员失败：%v", err)
	}

	rc.SRem(context.Background(), "group:"+gid, member)
	rc.SRem(context.Background(), "inGroup:"+member, gid)

	return nil
}

func kickAll(gid string) {
	groupMap, _ := ListGroupMember(gid)

	for uid := range groupMap {
		rc.SRem(context.Background(), "group:"+gid, uid)
		rc.SRem(context.Background(), "inGroup:"+uid, gid)
	}
}

func GetGroupReq(uid, gid string) (map[string]string, error) {
	if !IsGroupAdmin(uid, gid) && !IsGroupOwner(uid, gid) {
		return nil, fmt.Errorf("无权限")
	}
	groupReqList, err := rc.HGetAll(context.Background(), "groupReq:"+gid).Result()
	if err != nil {
		return nil, fmt.Errorf("获取群请求列表失败:%s", err)
	}

	return groupReqList, err
}

func RespGroupReq(uid, gid, isAccept string) error {
	_, err := rc.HGet(context.Background(), "groupReq:"+uid, gid).Result()
	if err == redis.Nil {
		return fmt.Errorf("不存在这个群申请")
	}

	if isAccept == "1" {
		rc.SAdd(context.Background(), "group:"+gid, uid)
		rc.SAdd(context.Background(), "inGroup:"+uid, gid)
	}

	rc.HDel(context.Background(), "groupReq:"+uid, gid)
	return nil
}

func ListGroupAdmin(gid string) (map[string]string, error) {
	memberList, err := rc.SMembers(context.Background(), "groupAdmin:"+gid).Result()
	if err != nil {
		return nil, fmt.Errorf("获取群成员列表失败:%s", err)
	}

	groupMap := make(map[string]string)

	for _, groupgid := range memberList {
		groupMap[groupgid] = MustGetNameById(groupgid)
	}
	return groupMap, err
}

func IsGroupAdmin(uid string, gid string) bool {
	// 查询发送者的群列表
	groupMap, _ := ListGroupAdmin(uid)

	// 判断接收者是否在群列表中
	return groupMap[gid] != ""
}

func IsGroupOwner(uid string, gid string) bool {
	group, err := GetGroupByGid(gid)
	if err != nil {
		return false
	}
	return group.Owner == uid
}
