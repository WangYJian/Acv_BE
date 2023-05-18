package api

import (
	"acv/processQueue"
	"acv/webSocket"
	"fmt"
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

	// 将文件放入队列
	processQueue.ProcessQueue.PushBack(params)
	fmt.Println(processQueue.ProcessQueue.GetQueue())
	c.String(200, "File %s uploaded successfully.", params.File.Filename)
}
