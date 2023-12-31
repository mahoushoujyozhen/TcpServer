package znet

import (
	"TcpServer/utils"
	"TcpServer/ziface"
	"errors"
	"fmt"
	"net"
	"time"
)

// IServer 接口实现，定义一个Server服务类
type Server struct {
	// 服务器的名称
	Name string
	// tcp4 or orther
	IPVersion string
	// 服务绑定的IP地址
	IP string
	// 服务绑定的端口
	Port int
	//当前Server由用户绑定的回调router，也就是Server注册的链接对应的处理业务
	MsgHandler ziface.IMsgHandle
	//当前Server的连接管理器
	ConnMgr ziface.IConnManager

	//新增两个hook函数原型,用来接收hook函数
	//该Server的链接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的链接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

// ============== 定义当前客户端链接的handle api ===========
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

//  =================实现ziface.IServer所有的接口=================

// 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	//开启一个go去做服务端Linster业务
	go func() {
		//0、启动worker工作池机制
		s.MsgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		//已经监听成功
		fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		//3 启动server网络连接业务
		for {
			/*
				思考一下，这里其实只有一个server，但是为什么在connection里面需要记录server这个句柄，只有一个server，
				需要用到server的时候，在这里执行就可以了，有必要定义一个句柄，让connection和server建立起互相索引的关系吗？
				目前我看着只能起一个server，而不是多个server
				答：这种设计可以满足多个server或者单个server的情况，通配性更强
			*/

			//3.1 阻塞等待客户端建立连接请求  这里监听新建的链接，同属于一个server
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//如果当前的连接数量已经达到上限，直接拒绝连接
				conn.Close()
				continue
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	//  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理

	//在server停止的时候，将全部的链接清除
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	// 阻塞，否则GO退出
	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// 创建一个服务器句柄，服务器句柄用来控制一个服务的连接，断开，就是一个服务的管理器
func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()
	//先初始化全局配置文件
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),   //msgHandler初始化
		ConnMgr:    NewConnManager(), //创建ConnManager
	}
	return s
}

// 根据msgId给服务添加路由
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router success! msgId = ", msgId)
}

// 得到连接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 设置该Server的链接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 设置该Server的链接断开时的hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用链接OnConnStart函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("=====> CallOnConnStart......")
		s.OnConnStart(conn)
	}
}

// 调用链接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("=====> CallOnConnStop......")
		s.OnConnStop(conn)
	}
}
