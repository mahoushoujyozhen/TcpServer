package main

import (
	"TcpServer/ziface"
	"TcpServer/znet"
	"fmt"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	//回写数据 这里reqeust利用句柄获取connection，这里如果request多的话，会造成短暂阻塞，所以需要有缓冲的读写channel
	err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// HelloZinxRouter Handle
type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端数据，在回写ping...ping...ping
	fmt.Println("reccv from client : msgId=", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router v0.8"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ...")

	//===========设置两个链接属性，在创建链接之后
	fmt.Println("Set conn Name, Home done !")
	conn.SetProperty("Name", "mahoushoujyo")
	conn.SetProperty("Age", 18)
	//================

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

//链接断卡ID时候执行

func DoConnectionLost(conn ziface.IConnection) {
	//===========在链接销毁之前，查询conn的Name，Age属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if age, err := conn.GetProperty("Age"); err == nil {
		fmt.Println("Conn Property Age = ", age)
	}
	//===========

	fmt.Println("DoConnectionLost is Called ...")
}

func main() {

	//这里可以用goroutine开启多个server服务

	//创建一个server句柄
	s := znet.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//开启服务
	s.Serve()
}
