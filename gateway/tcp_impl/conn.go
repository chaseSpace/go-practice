package tcp_impl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"service/gateway2/gutil"
	"service/gateway2/tcp_iface"
	"sync"
)

type Connection struct {
	ctx context.Context
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID string
	//当前连接的关闭状态
	isClosed bool
	//该连接的处理方法router
	Router tcp_iface.IRouter
	//告知该链接已经退出/停止的channel
	ExitChan chan bool
	//消息处理
	h       tcp_iface.IMsgHandle
	msgChan chan []byte
	// conn关闭时通知server
	serverChan chan tcp_iface.IConnection

	closeNow bool

	onConnStart func(tcp_iface.IConnection)
	onConnStop  func(tcp_iface.IConnection)

	vmap sync.Map
}

func NewConnection(conn *net.TCPConn, connID string, h tcp_iface.IMsgHandle, serverChan chan tcp_iface.IConnection, opt ...ConnOpt) tcp_iface.IConnection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		h:          h,
		ExitChan:   make(chan bool),
		msgChan:    make(chan []byte, 10), // 允许每个conn缓冲10条消息，并保证在发生错误时将消息尽量发送出去
		serverChan: serverChan,
		ctx:        context.TODO(),
	}
	for _, o := range opt {
		o(c)
	}
	return c
}

type ConnOpt func(*Connection)

func WithCloseNow(b bool) func(*Connection) {
	return func(c *Connection) {
		c.closeNow = b
	}
}

func WithOnConnStart(f func(tcp_iface.IConnection)) func(c *Connection) {
	return func(c *Connection) {
		c.onConnStart = f
	}
}

func WithOnConnStop(f func(tcp_iface.IConnection)) func(c *Connection) {
	return func(c *Connection) {
		c.onConnStop = f
	}
}

func (c *Connection) StartReader() {
	fmt.Printf("new conn-%s Reader Goroutine is  running\n", c.ConnID)

	for {
		// 创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head  TODO 使用pool优化
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetConn(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		//拆包，得到msgid 和 datalen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		//根据 dataLen 读取 data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			// TODO 使用pool优化
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetConn(), data); err != nil {
				fmt.Println("read msg msg error ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			conn:     c,
			IMessage: msg, //将之前的buf 改成 msg
		}

		if gutil.Conf.Tcp.WorkerPoolSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.h.MsgEnqueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.h.HandleMsg(&req)
		}
	}

	close(c.ExitChan)
}

func (c *Connection) StartWriter() {
	for {
		select {
		case msg := <-c.msgChan:
			_, err := c.Conn.Write(msg)
			if err != nil {
				fmt.Printf("conn-%s write err:%v\n", c.ConnID, err)
				continue
			}
			//fmt.Printf("conn-%s write ok\n", c.ConnID)
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	defer c.Stop()

	if c.closeNow {
		_ = c.SendMsg(0, []byte("server is busy!"), false)
		return
	}
	// 一个conn做读写分离，内部使用带缓冲chan通信
	go c.StartReader()
	go c.StartWriter()

	c.onConnStart(c)

	<-c.ExitChan
	//for {
	//	select {
	//	case <-c.ExitChan:
	//		//得到退出消息，不再阻塞
	//		return
	//	}
	//}
}

func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	defer c.onConnStop(c)

	c.isClosed = true
	// 发完所有数据 再关闭conn
	restDataLen := len(c.msgChan)
	for i := 0; i < restDataLen; i++ {
		msg := <-c.msgChan
		_, err := c.Conn.Write(msg)
		if err != nil {
			fmt.Printf("Close: conn-%s write err:%v\n", c.ConnID, err)
		}
	}
	//TODO Connection Close() 如果用户注册了该链接的关闭回调业务，那么在此刻应该调用
	// 关闭socket链接
	_ = c.Conn.Close()

	fmt.Println("conn", c.ConnID, "exited!")

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.serverChan <- c
}

func (c *Connection) GetConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() string {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte, byChan bool) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data封包，然后传给chan，等待发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	if !byChan {
		_, err := c.Conn.Write(msg)
		if err != nil {
			fmt.Printf("byChan conn-%s write err:%v\n", c.ConnID, err)
			return err
		}
	} else {
		c.msgChan <- msg
	}
	return nil
}

func (c *Connection) SetKey(key string, val any) {
	c.vmap.Store(key, val)
}

func (c *Connection) GetKey(key string) any {
	v, _ := c.vmap.Load(key)
	return v
}
