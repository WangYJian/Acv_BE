package api

import (
	"acv/util"
	"acv/webSocket"
	"log"

	"github.com/gin-gonic/gin"
)

// Id id结构体
type Id struct {
	Id string `json:"id"`
}

// Connect 连接
func Connect(c *gin.Context) {
	// 升级协议
	conn, err := webSocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error occurred while upgrading protocol: %s", err.Error())
		return
	}

	// 生成id
	id := util.GenerateLinkId()

	// 将id和连接放入map
	webSocket.IdConnMap[id] = conn

	// 返回id
	Id := Id{
		Id: id,
	}
	err = conn.WriteJSON(Id)
	if err != nil {
		log.Printf("Error occurred while writing message: %s", err.Error())
	}
}
