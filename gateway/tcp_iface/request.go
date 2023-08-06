package tcp_iface

import "context"

type IRequest interface {
	IMessage                    //获取请求消息的数据
	GetConnection() IConnection //获取请求连接信息
	Ctx() context.Context
}
