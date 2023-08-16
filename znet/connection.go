package znet

import (
	"TcpServer/utils"
	"TcpServer/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	//句柄，当前Conn属于哪个Server  ，有时候在conn中，我们需要server.ConnMgr的使用权，所以需要知道属于哪个Server
	TcpServer ziface.IServer //当前conn属于哪个server，在connection初始化的时候添加即可
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID，也可以成为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandle
	//告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
	//无缓冲通道，用于读，写两个goroutine之间的消息通信  进行读写分离设计
	msgChan chan []byte
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
}

// 创建连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1), //用来通知主进程是否能结束，防止主进程结束然后goroutine意外死亡
		msgChan:      make(chan []byte),  //msgChan初始化
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}
	//在创建连接的时候，将conn添加到链接管理中，
	//这里需要用到ConnMgr，所以体现了TcpServer句柄的作用,这个句柄还可能在其他场景下发挥作用
	c.TcpServer.GetConnMgr().Add(c) //将当前新建的链接添加到ConnMgr中
	return c
}

/*
	写消息Goroutine，server将用户数据发送给客户端
*/

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给前端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error: ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitBuffChan:
			//conn已经关闭
			return
		}
	}
}

// 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader GoRoutine is running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " conn reader exit")
	//注意，这个版本的话不会退出for循环，所以这里的Stop不会执行，需要考虑什么情况下会退出for循环
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经启动工作池机制，将request发送给工作池处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	//开启用户从客户端读取数据流程的Gorouitne
	go c.StartReader()
	//开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

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

	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	//1、如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true //设置关闭连接的标识位
	//TODO Connection Stop() 如果用户注册了该连接的关闭回调业务，那么在此刻应该显式调用

	//关闭socket连接
	c.Conn.Close()

	//关闭writer Goroutine
	c.ExitBuffChan <- true

	//将链接从链接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c) //从ConnManager中删除conn

	//关闭该连接的所有管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
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

//直接将Message数据发送到msgChan

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
	c.msgChan <- msg //这里写到msgChan促发Write Goroutine将数据写会客户端
	return nil
}
