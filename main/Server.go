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

	//回写数据
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
	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router v0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	//这里可以用goroutine开启多个server服务

	//创建一个server句柄
	s := znet.NewServer()

	//配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//开启服务
	s.Serve()
}
