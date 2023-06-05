package api

import (
	"acv/client/db"
	"acv/model"
	"acv/util"
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	UserID          string `json:"userID" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// Login 登录 API
func Login(c *gin.Context) {
	// 解析请求体中的 JSON 数据到 LoginRequest 结构体中
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		log.Println("Error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据用户账户查找用户
	result, err := model.FindUser(db.DbEngine, loginReq.UserID)
	if err != nil {
		log.Println("Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil {
		log.Println("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 验证密码是否正确
	md5Password := md5.Sum([]byte(loginReq.Password))
	if hex.EncodeToString(md5Password[:]) != result.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不正确"})
		return
	}

	// 生成 token 并返回
	token := util.GenerateToken(loginReq.UserID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Register 注册 API
func Register(c *gin.Context) {
	// 解析请求体中的 JSON 数据到 RegisterRequest 结构体中
	var registerReq RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查账号是否重复
	result, _ := model.FindUser(db.DbEngine, registerReq.UserID)
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "账号已存在"})
		return
	}

	// 检查密码和确认密码是否一致
	if registerReq.Password != registerReq.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码和确认密码不一致"})
		return
	}

	// 将账号和加密后的密码存入数据库
	model.AddUser(db.DbEngine, registerReq.UserID, registerReq.Password)

	// 生成并返回用户 token
	token := util.GenerateToken(registerReq.UserID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
