package db

import (
	"c6m/model"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func CreateGroup(uid, gname string) (*model.Group, error) {
	owner, _ := GetGidByGname(gname)
	if owner != "" {
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

	return nil
}

func GetGroupReq(gid string) (map[string]string, error) {
	groupReqList, err := rc.HGetAll(context.Background(), "groupReq:"+gid).Result()
	if err != nil {
		return nil, fmt.Errorf("获取群请求列表失败:%s", err)
	}

	return groupReqList, err
}

func RespGroupReq(gid string, groupgid string, isAccept string) error {
	_, err := rc.HGet(context.Background(), "groupReq:"+gid, groupgid).Result()
	if err == redis.Nil {
		return fmt.Errorf("不存在这个群申请")
	}

	if isAccept == "1" {
		rc.SAdd(context.Background(), "group:"+gid, groupgid)
		rc.SAdd(context.Background(), "group:"+groupgid, gid)
	}

	rc.HDel(context.Background(), "groupReq:"+gid, groupgid)
	return nil
}

func setAdmin() {}
