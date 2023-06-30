package tcp_impl

import (
	"code.com/tars/goframework/kissgo/appzaplog"
	"code.com/tars/goframework/kissgo/appzaplog/zap"
	"fmt"
	"net"
	"service/gateway2/gutil"
	"service/gateway2/tcp_iface"
	"time"
	"util/common"
)

type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//消息处理
	h tcp_iface.IMsgHandle

	connClosedChan chan tcp_iface.IConnection

	//连接管理器
	ConnMgr tcp_iface.IConnManager
}

func NewServer(name string) tcp_iface.IServer {
	s := &Server{
		Name:           name + "-" + string(common.RandomBytes2(0, 3)),
		IPVersion:      "tcp4",
		IP:             "0.0.0.0",
		Port:           gutil.Conf.Tcp.Port,
		h:              NewMsgHandle(),
		connClosedChan: make(chan tcp_iface.IConnection),
		ConnMgr:        NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	hint := "ServerStart-"
	go func() {
		s.h.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		//已经监听成功
		fmt.Println("start Tcp server  ", s.Name, " succ, now listening...")

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			if err = setConn(conn, hint); err != nil {
				continue
			}

			connId := gutil.GenConnId()

			closeNow := s.ConnMgr.Len() >= gutil.Conf.Tcp.MaxConn

			var opts []ConnOpt
			opts = append(opts, WithCloseNow(closeNow))
			opts = append(opts, WithOnConnStart(onConnStart))
			opts = append(opts, WithOnConnStop(onConnStop))

			// 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(conn, connId, s.h, s.connClosedChan, opts...)
			s.ConnMgr.Add(dealConn)

			// 启动当前链接的处理业务
			go dealConn.Start()

			fmt.Println("connection counter=", "closeNow=", closeNow)
		}
	}()
}

func (s *Server) AddRouter(msgId uint32, router tcp_iface.IRouter) {
	s.h.AddRouter(msgId, router)

	fmt.Println("Add Router msgID=", msgId, "ok!")
}

// TODO 有用？
func (s *Server) GetConnMgr() tcp_iface.IConnManager {
	return s.ConnMgr
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Tcp server , name ", s.Name)

	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	defer s.Stop()

	for {
		select {
		case conn := <-s.connClosedChan:
			s.ConnMgr.Remove(conn)
			println("Server connID exited:", conn.GetConnID())
		}
	}
}

// ------------------------------

func onConnStart(c tcp_iface.IConnection) {
	fmt.Println("onConnStart")
}

func onConnStop(tcp_iface.IConnection) {
	fmt.Println("onConnStop")
}

func setConn(conn *net.TCPConn, hint string) error {
	err := conn.SetLinger(3)
	if err != nil {
		appzaplog.Error(hint+"setConn SetLinger", zap.Error(err), zap.String("addr", conn.RemoteAddr().String()))
		return err
	}
	err = conn.SetKeepAlive(true)
	if err != nil {
		appzaplog.Error(hint+"setConn SetKeepAlive", zap.Error(err), zap.String("addr", conn.RemoteAddr().String()))
		return err
	}
	err = conn.SetKeepAlivePeriod(time.Second * 20)
	if err != nil {
		appzaplog.Error(hint+"setConn SetKeepAlivePeriod", zap.Error(err), zap.String("addr", conn.RemoteAddr().String()))
		return err
	}
	// 使用系统默认值
	//err = conn.SetReadBuffer(2000)
	//err = conn.SetWriteBuffer(2000)
	return nil
}
