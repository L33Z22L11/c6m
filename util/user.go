package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

func validateTokenAndGetUID(token string) (string, error) {
	// 检查令牌是否为空
	if token == "" {
		return "", fmt.Errorf("未提供访问令牌")
	}

	// 从请求头部的Authorization字段中提取令牌
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", fmt.Errorf("无效的访问令牌")
	}
	token = tokenParts[1]

	// 从Redis中根据访问令牌获取用户ID
	uid, err := rc.Get(context.Background(), token).Result()
	if err != nil {
		if err == redis.Nil {
			// 令牌不存在或已过期
			return "", fmt.Errorf("无效的访问令牌")
		}
		return "", fmt.Errorf("从数据库获取用户ID失败: %v", err)
	}

	return uid, nil
}
