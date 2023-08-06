package tcp_iface

import (
	"net"
)

// 定义连接接口
type IConnection interface {
	Start()
	Stop()
	GetConn() *net.TCPConn
	GetConnID() string
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte, byChan bool) error
	SetKey(string, any)
	GetKey(string) any
}
