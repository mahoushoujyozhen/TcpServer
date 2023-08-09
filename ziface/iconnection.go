package ziface

import "net"

type IConnection interface {
	//启动连接
	Start()
	//停止连接，结束当前的连接状态
	Stop()
	//从当前连接获取原始的socket TCPConn GetTCOConnection() *net.TCPConn
	// 获取当前连接的ID
	GetConnID() uint32
}

// 定义一个统一处理链路业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
