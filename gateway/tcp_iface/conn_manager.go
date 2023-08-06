package tcp_iface

/*
IConnManager
连接管理抽象层
*/
type IConnManager interface {
	Add(conn IConnection)                   //添加链接
	Remove(conn IConnection)                //删除连接
	Get(connID string) (IConnection, error) //利用ConnID获取链接
	Len() int                               //获取当前连接数
	ClearConn()                             //删除并停止所有链接
}
