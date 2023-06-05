package api

import (
	"acv/client/db"
	"acv/model"
	"acv/processQueue"
	"acv/webSocket"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Upload 上传文件
func Upload(c *gin.Context) {
	// 获取参数
	var params processQueue.Job
	_ = c.ShouldBind(&params)

	// 查找id对应的连接
	conn := webSocket.IdConnMap[params.Id]
	// 如果连接不存在
	if conn == nil {
		c.String(400, "Connection does not exist.")
		return
	}
	// 如果连接存在，将连接放入参数
	params.Conn = conn

	// 修改文件名，加上当前时间戳
	params.File.Filename = fmt.Sprintf("%d_%s", time.Now().Unix(), params.File.Filename)

	// 将文件保存到本地
	path := os.Getenv("FILE_PATH")
	err := c.SaveUploadedFile(params.File, path+"/"+params.File.Filename)
	if err != nil {
		c.String(400, "File upload failed.")
		return
	}

	// 将信息保存到数据库
	file := model.File{
		LinkID:     params.Id,
		UserID:     c.GetString("id"),
		Filename:   params.File.Filename,
		Path:       path,
		CreateTime: time.Now(),
	}

	// 将文件信息保存到数据库
	err = model.CreateFile(db.DbEngine, file)
	if err != nil {
		c.String(400, "File upload failed.")
		return
	}

	// 将文件放入队列
	processQueue.ProcessQueue.PushBack(params)
	fmt.Println(processQueue.ProcessQueue.GetQueue())
	c.String(200, "%s 文件上传成功", params.File.Filename)
}
