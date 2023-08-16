package znet

import (
	"TcpServer/ziface"
	"errors"
	"fmt"
	"sync"
)

/*
	连接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理连接，用connID来管理
	connLock    sync.RWMutex                  //读写连接的读写锁  map读写不是并发安全的，需要读写锁来控制保证安全性
}

/*
	创建一个连接管理
*/

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//添加连接

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("[connection id] ", conn.GetConnID(), " add to ConnManager successfully : conn num = ",
		connMgr.Len())
}

// 删除连接 单纯的从Map中移除，但是conn没有停止
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息  根据map 的key 去删除对应的连接
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("[connection id] ", conn.GetConnID(), " Remove successfully : conn num = ",
		connMgr.Len())
}

// 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源Map  加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

// 获取当前连接数量
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除所有的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}
