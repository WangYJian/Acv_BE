package util

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func GenerateLinkId() string {
	currentTime := time.Now().String()
	hash := md5.Sum([]byte(currentTime))
	return hex.EncodeToString(hash[:])
}
