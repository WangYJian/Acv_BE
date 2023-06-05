package processQueue

import (
	"acv/client/db"
	"acv/model"
	"acv/webSocket"
	"container/list"
	"fmt"
	"mime/multipart"
	"strconv"
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

var ProcessQueue *Queue

// InitQueue 初始化任务队列
func InitQueue() {
	ProcessQueue = &Queue{
		queueList: list.New(),
	}
	go ProcessQueue.Process()
}

// PushBack 将任务放入队列
func (q *Queue) PushBack(job Job) {
	q.queueList.PushBack(job)
	// 通知用户在队列中的位置
	NotifyUserByOrder(q.queueList.Len(), job)
}

// Process 处理队列中的任务
func (q *Queue) Process() {
	for {
		if q.queueList.Len() > 0 {
			// 取出任务
			job := q.queueList.Front().Value.(Job)
			// 处理任务
			result := processFile(job.File)
			// 从数据库获取文件信息
			file, _ := model.GetFileByLinkID(db.DbEngine, job.Id)
			// 更新文件信息
			file.Result = result
			file.Update(db.DbEngine)
			// 通知用户
			notifyUser(job, job.File.Filename+result)
			// 删除任务
			q.queueList.Remove(q.queueList.Front())
		}
	}
}

// GetQueue 按照队列顺序获取连接
func (q *Queue) GetQueue() []Job {
	var connList []Job
	for e := q.queueList.Front(); e != nil; e = e.Next() {
		connList = append(connList, e.Value.(Job))
	}
	return connList
}

// processFile 处理文件
func processFile(file *multipart.FileHeader) string {
	time.Sleep(5 * time.Second)
	fmt.Println(file.Filename + "文件处理完毕")
	return "no error"

}

// notifyUser 通知用户
func notifyUser(job Job, message string) {
	// 从map中取出连接
	if job.Conn != nil {
		// 存在则发送消息
		msgStruct := msg{
			Id:  job.Id,
			Msg: message,
		}
		_ = job.Conn.WriteJSON(msgStruct)
		webSocket.Lock.Lock()
		// 通知用户关闭这个连接
		_ = job.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		_ = job.Conn.Close()
		// 从map中删除这个连接
		delete(webSocket.IdConnMap, job.Id)
		webSocket.Lock.Unlock()
	}
	// 遍历队列，给其他所有连接发送消息，告知它们的排队位置
	for i, iJob := range ProcessQueue.GetQueue() {
		if i == 0 {
			continue
		} else {
			NotifyUserByOrder(i, iJob)
		}
	}
}

// NotifyUserByOrder 按照队列顺序通知用户
func NotifyUserByOrder(i int, iJob Job) {
	var content string
	if i == 0 {
		return
	} else if i == 1 {
		content = "文件正在处理中"
	} else {
		content = "排在第" + strconv.Itoa(i) + "位"
	}
	msgStruct := msg{
		Id:  iJob.Id,
		Msg: iJob.File.Filename + content,
	}
	_ = iJob.Conn.WriteJSON(msgStruct)
}
