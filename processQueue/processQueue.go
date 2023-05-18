package processQueue

import (
	"acv/webSocket"
	"container/list"
	"fmt"
	"github.com/gorilla/websocket"
	"mime/multipart"
	"strconv"
	"time"
)

// message 消息结构体
type msg struct {
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
}

// Process 处理队列中的任务
func (q *Queue) Process() {
	for {
		if q.queueList.Len() > 0 {
			// 取出任务
			job := q.queueList.Front().Value.(Job)
			// 处理任务
			processFile(job.File)
			// 通知用户
			notifyUser(job, "Your file has been processed.")
			// 删除任务
			q.queueList.Remove(q.queueList.Front())
		}
	}
}

// GetQueue 按照队列顺序获取连接
func (q *Queue) GetQueue() []*websocket.Conn {
	var connList []*websocket.Conn
	for e := q.queueList.Front(); e != nil; e = e.Next() {
		connList = append(connList, e.Value.(Job).Conn)
	}
	return connList
}

// processFile 处理文件
func processFile(file *multipart.FileHeader) {
	time.Sleep(5 * time.Second)
	fmt.Println("File " + file.Filename + " processed.")
}

// notifyUser 通知用户
func notifyUser(job Job, message string) {
	// 从map中取出连接
	if job.Conn != nil {
		// 存在则发送消息
		msgStruct := msg{
			Msg: message,
		}
		_ = job.Conn.WriteJSON(msgStruct)
		webSocket.Lock.Lock()
		// 关闭这个连接
		_ = job.Conn.Close()
		// 从map中删除这个连接
		delete(webSocket.IdConnMap, job.Id)
		webSocket.Lock.Unlock()
	}
	// 遍历队列，给其他所有连接发送消息，告知它们的排队位置
	for i, conn := range ProcessQueue.GetQueue() {
		msgStruct := msg{
			Msg: "You are the " + strconv.Itoa(i) + "th in the queue.",
		}
		_ = conn.WriteJSON(msgStruct)
	}
}
