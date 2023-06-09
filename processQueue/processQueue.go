package processQueue

import (
	"acv/client/db"
	"acv/model"
	"acv/webSocket"
	"container/list"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gorilla/websocket"
)

// message 消息结构体
type msg struct {
	Id  string `json:"id"`
	Msg string `json:"message"`
}

type Job struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Conn *websocket.Conn
	Id   string `form:"id" binding:"required"`
}

type Queue struct {
	queueList *list.List
}

type Data struct {
	LinkID   string `json:"link_id"`
	Conn     *websocket.Conn
	Filename string `json:"filename"`
}

type Result struct {
	LinkID string `json:"link_id"`
	Conn   *websocket.Conn
	isFake bool `json:"is_fake"`
}

var ProcessQueue *Queue

// InitQueue 初始化任务队列
func InitQueue() {
	ProcessQueue = &Queue{
		queueList: list.New(),
	}
	go ProcessQueue.Process()
}

// PushBack 将任务放入队列
func (q *Queue) PushBack(data Data) {
	var str string
	if q.queueList.Len() == 0 {
		str = "正在处理中..."
	} else {
		str = "正在排队中..."
	}
	q.queueList.PushBack(data)
	// 通知用户正在排队
	msg := msg{
		Id:  data.LinkID,
		Msg: str,
	}
	// 将消息发送给用户
	_ = data.Conn.WriteJSON(msg)
}

// Process 处理队列中的任务
func (q *Queue) Process() {
	for {
		if q.queueList.Len() > 0 {
			// 将目前所有任务弹出，放入数组
			var data []Data
			for q.queueList.Len() > 0 {
				data = append(data, q.queueList.Front().Value.(Data))
				q.queueList.Remove(q.queueList.Front())
			}
			// 处理任务
			result := processFile(data)
			// 从数据库获取文件信息
			for _, v := range result {
				file, err := model.GetFileByLinkID(db.DbEngine, v.LinkID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				// 更新文件信息
				var str string
				if v.isFake == true {
					str = "有风险"
				} else {
					str = "正常"
				}
				file.Result = str
				file.Update(db.DbEngine)
			}
			// 通知用户
			notifyUser(result)
		}
	}
}

// GetQueue 按照队列顺序获取连接
func (q *Queue) GetQueue() []Data {
	var connList []Data
	for e := q.queueList.Front(); e != nil; e = e.Next() {
		connList = append(connList, e.Value.(Data))
	}
	return connList
}

// processFile 处理文件
func processFile(data []Data) []Result {
	time.Sleep(20 * time.Second)
	// 模拟处理文件
	var result []Result
	for _, v := range data {
		// 模拟结果
		result = append(result, Result{
			LinkID: v.LinkID,
			Conn:   v.Conn,
			isFake: true,
		})
	}
	return result
}

// notifyUser 通知用户
func notifyUser(result []Result) {
	fmt.Println(result)
	// 对每个结果进行处理
	for _, v := range result {
		// 从map中取出连接
		if v.Conn != nil {
			// 存在则发送消息
			var str string
			if v.isFake == true {
				str = "有风险"
			} else {
				str = "正常"
			}
			msgStruct := msg{
				Id:  v.LinkID,
				Msg: str,
			}
			_ = v.Conn.WriteJSON(msgStruct)
			webSocket.Lock.Lock()
			// 通知用户关闭这个连接
			_ = v.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			_ = v.Conn.Close()
			// 从map中删除这个连接
			delete(webSocket.IdConnMap, v.LinkID)
			webSocket.Lock.Unlock()
		}
	}

	// 通知其他用户正在处理
	for _, v := range webSocket.IdConnMap {
		msgStruct := msg{
			Id:  "all",
			Msg: "正在处理中...",
		}
		_ = v.WriteJSON(msgStruct)
	}
}
