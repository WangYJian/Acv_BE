package router

import (
	"acv/middleware"
	"acv/router/api"

	"github.com/gin-gonic/gin"
)

func UseMyRouter(r *gin.Engine) {
	acvApi := r.Group("/api")
	{
		// 登录
		acvApi.POST("/login", api.Login)
		// 注册
		acvApi.POST("/register", api.Register)
		// 上传
		acvApi.POST("/upload",
			middleware.Auth(),
			api.Upload,
		)
		// 连接
		acvApi.GET("/connect",
			middleware.Auth(),
			api.Connect,
		)
		// 查看历史记录
		acvApi.GET("/history",
			middleware.Auth(),
			api.History,
		)
		// 获取音频文件
		acvApi.GET("/audio",
			middleware.Auth(),
			api.GetAudioFile,
		)
	}
}
