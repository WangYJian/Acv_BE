package util

import (
	"encoding/base64"
	"strconv"
	"time"
)

func GenerateToken(userID string) string {
	// 获取当前时间戳
	timestamp := time.Now().Unix()

	// 将 userID 和时间戳拼接成字符串
	data := userID + strconv.FormatInt(timestamp, 10)

	// 对字符串进行 base64 编码
	token := base64.StdEncoding.EncodeToString([]byte(data))

	return token
}

func ParseToken(token string) (userID string, err error) {
	// 对 token 进行 base64 解码
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	// 从解码后的数据中解析出 userID 和时间戳
	userID = string(data[:len(data)-10])

	return userID, nil
}
