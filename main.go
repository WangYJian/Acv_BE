package main

import (
	"acv/processQueue"
	"acv/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

var port string

func SettingUpEnvironment() {
	// 读取配置文件
	godotenv.Load()
	// 配置端口
	port = os.Getenv("ACV_PORT")
	// 配置数据库
	//db.InitDB()
	// 初始化并运行任务队列
	processQueue.InitQueue()
}

func main() {
	// 初始化环境
	SettingUpEnvironment()
	// 初始化路由
	r := gin.Default()
	router.UseMyRouter(r)
	des := ":" + port
	_ = r.Run(des)
}
