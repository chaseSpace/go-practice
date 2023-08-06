package tcp_impl

import (
	"errors"
	"fmt"
	"service/gateway2/tcp_iface"
	"sync"
)

/*
ConnManager
连接管理模块
*/
type ConnManager struct {
	cmap  map[string]tcp_iface.IConnection //管理的连接信息
	mutex sync.RWMutex                     //读写连接的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		cmap: make(map[string]tcp_iface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn tcp_iface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.mutex.Lock()
	defer connMgr.mutex.Unlock()

	//将conn连接添加到ConnManager中
	connMgr.cmap[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len(), conn.GetConnID())
}

func (connMgr *ConnManager) Remove(conn tcp_iface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.mutex.Lock()
	defer connMgr.mutex.Unlock()

	conn.Stop()
	delete(connMgr.cmap, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID string) (tcp_iface.IConnection, error) {
	//保护共享资源Map 加读锁
	connMgr.mutex.RLock()
	defer connMgr.mutex.RUnlock()

	if conn, ok := connMgr.cmap[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.cmap)
}

func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.mutex.Lock()
	defer connMgr.mutex.Unlock()

	clears := connMgr.Len()
	//停止并删除全部的连接信息
	for connID, conn := range connMgr.cmap {
		conn.Stop()
		delete(connMgr.cmap, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", clears)
}
