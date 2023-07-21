package database

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-redis/redis/v8"
)

// 定义 Redis 客户端
var rc *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "qmLwjsTgxMXJ39w4WkRFbPrc9Cgxjbgb",
	DB:       0,
})

func StartServer() {
	cmd := exec.Command("redis-server", "conf/redis.conf")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Redis服务器启动失败: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Redis服务器,启动!")
}

func Test() {
	if _, err := rc.Ping(rc.Context()).Result(); err != nil {
		fmt.Printf("Redis测试失败: %s\n", err)
		return
	}
	fmt.Println("Redis测试成功")
}

func StopServer() {
	if _, err := rc.Shutdown(rc.Context()).Result(); err != nil {
		fmt.Printf("Redis服务器关闭失败: %s\n", err)
		return
	}
	fmt.Println("Redis服务器关闭")
}
