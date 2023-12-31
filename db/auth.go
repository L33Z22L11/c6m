package db

import (
	"c6m/model"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// 定义用户结构体
func CreateUser(username, password, question, answer string) (*model.Auth, error) {
	// 验证用户名和密码是否符合规范
	if !isValidName(username) {
		return nil, fmt.Errorf("用户名[%s]格式错误: 长度3~18, 仅限数字、字母和下划线", username)
	}

	owner, _ := GetUidByUname(username)
	if owner != "" {
		return nil, fmt.Errorf("用户已注册")
	}

	if !isValidPw(password) {
		return nil, fmt.Errorf("密码格式错误: 长度6~32, 仅限常见字符(ASCII 32~126),")
	}

	if len(answer) == 0 {
		return nil, fmt.Errorf("必须填写验证答案")
	}

	// 生成随机的 salt 值
	salt := generateSalt()

	// 创建新用户
	auth := &model.Auth{
		Uid:      generateUid(),
		Username: username,
		Hash:     generateHash(password, salt),
		Salt:     salt,
		Question: question,
		Answer:   answer,
	}

	// 将用户保存到 Redis 哈希表中
	err := SaveAuth(auth)
	if err != nil {
		return nil, fmt.Errorf("保存用户[%s]到Redis: %v", username, err)
	}

	return auth, nil
}

func isValidName(username string) bool {
	isVaildCharset := regexp.MustCompile("^[0-9A-Z_a-z]+$").MatchString(username)
	haveProperLength := len(username) >= 3 && len(username) <= 18
	return isVaildCharset && haveProperLength
}

func isValidPw(password string) bool {
	isVaildCharset := regexp.MustCompile("^[ -~]+$").MatchString(password)
	haveProperLength := len(password) >= 6 && len(password) <= 32
	return isVaildCharset && haveProperLength
}

// 生成随机的 salt 值
func generateSalt() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// 使用 salt 和密码生成哈希值
func generateHash(password, salt string) string {
	// 将密码和 salt 拼接起来
	saltedPassword := []byte(password + salt)
	// 生成哈希值
	hash := sha256.Sum256(saltedPassword)
	// 将哈希值转换为 16 进制字符串
	return hex.EncodeToString(hash[:])
}

func AuthUser(username, password string) (string, string, error) {
	// 根据用户名从数据库中获取用户信息
	uid, err := GetUidByUname(username)
	if err != nil {
		return "", "", fmt.Errorf("数据库中没有这个用户: %v", err)
	}
	auth, _ := GetAuthByUid(uid)

	// 使用用户的 salt 值和密码生成哈希值
	hash := generateHash(password, auth.Salt)

	// 比较生成的哈希值和用户存储的哈希值是否一致
	if hash != auth.Hash {
		return "", "", errors.New("密码错误")
	}

	// 生成访问令牌
	token, err := generateToken(auth.Uid)
	if err != nil {
		return "", "", fmt.Errorf("token 生成失败: %v", err)
	}

	return uid, token, nil
}

func ChangePassword(username, oldPassword, newPassword string) error {
	uid, err := GetUidByUname(username)
	if err != nil {
		return fmt.Errorf("用户[%s]不存在", username)
	}

	auth, _ := GetAuthByUid(uid)

	// 验证旧密码是否正确
	if generateHash(oldPassword, auth.Salt) != auth.Hash {
		return fmt.Errorf("旧密码不正确")
	}

	// 验证新密码是否符合规范
	if !isValidPw(newPassword) {
		return fmt.Errorf("新密码格式错误: 长度6~32, 仅限常见字符(ASCII 32~126)")
	}

	// 生成新的 salt 值
	newSalt := generateSalt()

	// 更新密码信息
	auth.Hash = generateHash(newPassword, newSalt)
	auth.Salt = newSalt

	// 更新用户信息到 Redis 哈希表中
	err = SaveAuth(auth)
	if err != nil {
		return fmt.Errorf("保存用户信息到Redis失败: %v", err)
	}

	return nil
}

func GetAuthQuestionByUname(username string) (string, error) {
	// 验证用户名是否存在
	uid, err := GetUidByUname(username)
	if err != nil {
		return "", fmt.Errorf("获取用户名失败: %v", err)
	}
	if uid == "" {
		return "", fmt.Errorf("用户名不存在")
	}

	// 获取用户的安全问题
	auth, _ := GetAuthByUid(uid)
	return auth.Question, nil
}

func ResetPwByAnswer(username, answer, newPassword string) error {
	// 验证用户名是否存在
	uid, err := GetUidByUname(username)
	if err != nil {
		return fmt.Errorf("获取用户名失败: %v", err)
	}
	if uid == "" {
		return fmt.Errorf("用户名不存在")
	}

	// 获取用户的安全问题答案
	auth, _ := GetAuthByUid(uid)
	if answer != auth.Answer {
		return fmt.Errorf("答案错误")
	}
	// 验证新密码是否符合规范
	if !isValidPw(newPassword) {
		return fmt.Errorf("新密码格式错误: 长度6~32, 仅限常见字符(ASCII 32~126)")
	}

	// 生成新的 salt 值
	newSalt := generateSalt()

	// 更新密码信息
	auth.Hash = generateHash(newPassword, newSalt)
	auth.Salt = newSalt

	// 更新用户信息到 Redis 哈希表中
	err = SaveAuth(auth)
	if err != nil {
		return fmt.Errorf("保存用户信息到Redis失败: %v", err)
	}

	return nil
}

// 生成 uid
func generateUid() string {
	uid := 1000
	lastUid, err := rc.Get(context.Background(), "lastUid").Result()
	if err != redis.Nil {
		uid, _ = strconv.Atoi(lastUid)
	}
	uid++
	rc.Set(context.Background(), "lastUid", uid, 0).Err()
	return fmt.Sprintf("u%d", uid)
}

func SaveAuth(auth *model.Auth) error {
	// 将用户信息序列化为 JSON 字符串
	authJson, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	// 将用户信息存储到 Redis 哈希表中
	err = rc.HSet(context.Background(), "uid", auth.Username, auth.Uid).Err()
	if err != nil {
		return err
	}
	err = rc.HSet(context.Background(), "auth", auth.Uid, authJson).Err()
	if err != nil {
		return err
	}

	return nil
}

// 生成访问令牌
func generateToken(uid string) (string, error) {
	// 生成一个新的 UUID 作为访问令牌
	token := uuid.New().String()

	// 将访问令牌存储到 Redis 中，以便后续验证
	err := rc.Set(context.Background(), token, uid, time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("向数据库存储 token 失败: %v", err)
	}

	return token, nil
}

func GetAuthByUid(uid string) (*model.Auth, error) {
	userJson, err := rc.HGet(context.Background(), "auth", uid).Bytes()
	if err != nil {
		return nil, err
	}

	// 将 JSON 字符串反序列化为用户结构体
	var auth model.Auth
	err = json.Unmarshal(userJson, &auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func GetUidByUname(username string) (string, error) {
	uid, err := rc.HGet(context.Background(), "uid", username).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("用户[%s]不存在", username)
	}
	return uid, nil
}

func GetUidByToken(token string) (string, error) {
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

// 生成 gid
func generateGid() string {
	gid := 1000
	lastGid, err := rc.Get(context.Background(), "lastGid").Result()
	if err != redis.Nil {
		gid, _ = strconv.Atoi(lastGid)
	}
	gid++
	rc.Set(context.Background(), "lastGid", gid, 0).Err()
	return fmt.Sprintf("g%d", gid)
}

func GetGroupByGid(gid string) (*model.Group, error) {
	groupJson, err := rc.HGet(context.Background(), "group", gid).Bytes()
	if err != nil {
		return nil, err
	}

	// 将 JSON 字符串反序列化为群结构体
	var group model.Group
	err = json.Unmarshal(groupJson, &group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func GetGidByGname(gname string) (string, error) {
	uid, err := rc.HGet(context.Background(), "gid", gname).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("群[%s]不存在", gname)
	}
	return uid, nil
}

func MustGetNameById(id string) string {
	switch id[0] {
	case 'u':
		auth, _ := GetAuthByUid(id)
		return auth.Username
	case 'g':
		group, _ := GetGroupByGid(id)
		return group.Gname
	}
	return "[InvalidName]"
}
