package database

import (
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

	isFriendAdded, isUserAdded := true, true
	_, err = rc.HGet(context.Background(), fmt.Sprintf("friend:%s", uid), friendUid).Result()
	if err == redis.Nil {
		isFriendAdded = false
	}
	_, err = rc.HGet(context.Background(), fmt.Sprintf("friend:%s", friendUid), uid).Result()
	if err == redis.Nil {
		isUserAdded = false
	}

	if isFriendAdded {
		return fmt.Errorf("不得重复发送好友申请")
	}

	if !isUserAdded {
		rc.HSet(context.Background(), fmt.Sprintf("friend:%s", uid), friendUid, "")
		rc.HSet(context.Background(), fmt.Sprintf("friendReq:%s", friendUid), uid, "")
		// SendFriendRequst()
	}

	return nil
}

func SendFriendRequst() {
	panic("unimplemented")
}
