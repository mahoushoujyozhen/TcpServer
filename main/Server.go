package main

import (
	"TcpServer/ziface"
	"TcpServer/znet"
	"fmt"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter // 一定要先基础BaseRouter
}

// Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle ")
	_, err := request.GetConnection().(*znet.Connection).GetTCPConnection().Write([]byte("before ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error! ")
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle ")
	//使用断言将接口类型转换为特定的结构体，然后访问结构体特有的元素，如果这里不用断言，用不了GetTCPConnection方法
	//GetConnection()方法返回的是接口类型，可以用断言转换为实现了该接口的结构体Connection
	//断言为*Connection类型，才可以调用对应结构体的方法
	_, err := request.GetConnection().(*znet.Connection).GetTCPConnection().Write([]byte("ping ...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error! ")
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().(*znet.Connection).GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

// Server 模块的测试函数
func main() {

	//1 创建一个server 句柄 s
	s := znet.NewServer("[zinx V0.1]")

	s.AddRouter(&PingRouter{})
	//2 开启服务
	s.Serve()
}
