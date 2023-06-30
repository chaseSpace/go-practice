package tcp_impl

import (
	"fmt"
	"service/gateway2/gutil"
	"service/gateway2/tcp_iface"
	"strconv"
	"time"
)

// MsgHandle 处理全局的TCP消息发送
type MsgHandle struct {
	Apis           map[uint32]tcp_iface.IRouter //存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize int64                        //业务工作Worker池的数量，控制worker的全局数量
	TaskQueue      []chan tcp_iface.IRequest    //Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]tcp_iface.IRouter),
		WorkerPoolSize: gutil.Conf.Tcp.WorkerPoolSize,
		TaskQueue:      make([]chan tcp_iface.IRequest, gutil.Conf.Tcp.WorkerPoolSize),
	}
}

func (mh *MsgHandle) HandleMsg(request tcp_iface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is not FOUND!")
		return
	}
	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router tcp_iface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan tcp_iface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.HandleMsg(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	println("StartWorkerPool WorkerPoolSize=?", mh.WorkerPoolSize, "MaxWorkerTaskLen", gutil.Conf.Tcp.MaxWorkerTaskLen)
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan tcp_iface.IRequest, gutil.Conf.Tcp.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// MsgEnqueue 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) MsgEnqueue(request tcp_iface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := time.Now().UnixNano() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgId(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}
