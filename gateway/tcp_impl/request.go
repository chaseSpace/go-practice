package tcp_impl

import (
	"golang.org/x/net/context"
	"service/gateway2/tcp_iface"
)

type Request struct {
	conn               tcp_iface.IConnection //已经和客户端建立好的 链接
	tcp_iface.IMessage                       //客户端请求的数据
	ctx                context.Context
}

func (r *Request) GetConnection() tcp_iface.IConnection {
	return r.conn
}

func (r *Request) Ctx() context.Context {
	return r.ctx
}
