package router

import (
	"acv/router/api"
	"github.com/gin-gonic/gin"
)

func UseMyRouter(r *gin.Engine) {
	acvApi := r.Group("/api")
	{
		// 上传
		acvApi.POST("/upload", api.Upload)
		// 连接
		acvApi.GET("/connect", api.Connect)
	}
}
