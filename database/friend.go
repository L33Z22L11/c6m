package database

import (
	"context"
	"fmt"
)

func AddFriend(username, friendName string) error {
	uid, err := GetUidByUname(username)
	if err != nil {
		return err
	}
	friendUid, err := GetUidByUname(friendName)
	if err != nil {
		return err
	}

	rc.SIsMember(context.Background(), fmt.Sprintf("friends:%s", uid), friendUid)
	rc.SIsMember(context.Background(), fmt.Sprintf("friends:%s", friendUid), uid)

	err = rc.SAdd(context.Background(), fmt.Sprintf("friends:%s", uid), friendUid).Err()

	return nil
}
