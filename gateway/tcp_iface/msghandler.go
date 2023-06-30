package tcp_iface

type IMsgHandle interface {
	HandleMsg(request IRequest)             //马上以非阻塞方式处理消息
	AddRouter(msgId uint32, router IRouter) //为消息添加具体的处理逻辑
	StartWorkerPool()                       //启动worker工作池
	MsgEnqueue(request IRequest)            //将消息写入队列,由worker进行处理
}
