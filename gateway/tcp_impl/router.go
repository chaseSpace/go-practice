package tcp_impl

import (
	"fmt"
	"service/gateway2/tcp_iface"
	"time"
)

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req tcp_iface.IRequest)  {}
func (br *BaseRouter) Handle(req tcp_iface.IRequest)     {}
func (br *BaseRouter) PostHandle(req tcp_iface.IRequest) {}

/* --------------------- 分割线 -------------------------*/

type PingRouter struct {
	BaseRouter
}

func (*PingRouter) Handle(request tcp_iface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from conn=", request.GetConnection().GetConnID(),
		" Data=", request.GetMsgId(), string(request.GetData()), time.Now().String())

	//回写数据
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("ping...ping...ping"), true)
	if err != nil {
		fmt.Println("PingRouter err:", err)
	}
}
