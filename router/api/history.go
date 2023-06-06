package api

import (
	"acv/client/db"
	"acv/model"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type response struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	LinkID   string `json:"linkID"`
}

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

	// 将files转换为result
	var result []response
	for _, file := range files {
		result = append(result, response{
			Title:    file.Filename,
			Subtitle: file.Result,
			LinkID:   file.LinkID,
		})
	}

	// 将文件历史作为 JSON 数据返回给前端
	c.JSON(http.StatusOK, gin.H{
		"result": result,
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

// DeleteFile 处理删除文件的请求。
func DeleteFile(c *gin.Context) {
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

	// 判断用户是否有权限删除该音频文件
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden",
		})
		return
	}

	// 从数据库中删除音频文件
	if err := model.DeleteFileByLinkID(db.DbEngine, linkID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
