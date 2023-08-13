package znet

import (
	"TcpServer/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID，也可以成为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//该连接的处理方法router
	Router ziface.IRouter
	//告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
}

// 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1), //用来通知主进程是否能结束，防止主进程结束然后goroutine意外死亡
	}
	return c
}

// 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader GoRoutine is running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " conn reader exit")
	defer c.Stop()
	for {
		//创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}
		//拆包，得到msgid 和 dataLen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err: ", err)
			c.ExitBuffChan <- true
			continue
		}

		//根据dataLen 读取data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error :", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		//设置 data
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg, //将之前的 buf 改为 msg
		}
		//从路由Routers中找到注册绑定Conn的对应Handle
		go func(request ziface.IRequest) {
			//执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//开启处理该连接读取到客户端数据之后到请求业务
	go c.StartReader()

	//这里阻塞主进程，等待goroutine完成任务，后者发出完成任务的信号，主进程才退出
	//为了防止主进程退出，然后其中开启的goroutine会直接断开的问题
	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出的消息，不再阻塞
			return
		}
	}
}

// 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1、如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true //设置关闭连接的标识位
	//TODO Connection Stop() 如果用户注册了该连接的关闭回调业务，那么在此刻应该显式调用

	//关闭socket连接
	c.Conn.Close()

	//通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true

	//关闭该连接的所有管道
	close(c.ExitBuffChan)
}

// 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//  获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//直接将Message数据发送数据给远程的TCP客户端

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	//写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn write error")
	}
	return nil
}
