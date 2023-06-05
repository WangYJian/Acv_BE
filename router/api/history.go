package api

import (
	"acv/client/db"
	"acv/model"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// History 处理获取文件历史的请求。
func History(c *gin.Context) {
	// 获取用户ID
	userID := c.GetString("id")

	// 从数据库中获取文件历史
	files, err := model.GetFilesByUserID(db.DbEngine, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 将文件历史作为 JSON 数据返回给前端
	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

// GetAudioFile 处理获取音频文件的请求。
func GetAudioFile(c *gin.Context) {
	// 获取链接 ID
	linkID := c.Query("linkID")

	// 获取用户 ID
	userID := c.GetString("id")

	// 从数据库中获取音频文件信息
	file, err := model.GetFileByLinkID(db.DbEngine, linkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 判断用户是否有权限获取该音频文件
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden",
		})
		return
	}

	// 获取音频文件的完整路径
	audioFilePath := filepath.Join(file.Path, file.Filename)

	// 返回音频文件给前端播放
	c.File(audioFilePath)
}
