package model

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint `gorm:"primaryKey"`
	UserID     string
	Password   string
	CreateTime time.Time
}

// FindUser 查找用户信息
func FindUser(db *gorm.DB, userID string) (*User, error) {
	var user User
	result := db.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// AddUser 添加用户信息
func AddUser(db *gorm.DB, userID string, password string) error {
	md5Password := md5.Sum([]byte(password))
	user := User{UserID: userID, Password: hex.EncodeToString(md5Password[:]), CreateTime: time.Now()}
	result := db.Create(&user)
	return result.Error
}
