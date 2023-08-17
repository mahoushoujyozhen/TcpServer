package ziface

import "net"

type IConnection interface {
	//启动连接
	Start()
	//停止连接，结束当前的连接状态
	Stop()
	//从当前连接获取原始的socket TCPConn GetTCOConnection() *net.TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接的ID
	GetConnID() uint32
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
	//直接将Message数据发送数据给远程的TCP客户端（无缓冲）
	SendMsg(msgId uint32, data []byte) error
	//直接将Message数据发送给远程的TCP客户端（有缓冲）
	//因为一个connection会有多个request，如果多个request使用无缓冲的读写通知channel，会导致短暂的阻塞，有缓冲用来解决这个问题
	SendBuffMsg(msgId uint32, data []byte) error

	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

// 定义一个统一处理链路业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
