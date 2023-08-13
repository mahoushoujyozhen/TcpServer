package main

import (
	"TcpServer/znet"
	"fmt"
	"io"
	"net"
)

// // ping test 自定义路由
// type PingRouter struct {
// 	znet.BaseRouter // 一定要先基础BaseRouter
// }

// // Test PreHandle
// func (this *PingRouter) PreHandle(request ziface.IRequest) {
// 	fmt.Println("Call Router PreHandle ")
// 	_, err := request.GetConnection().(*znet.Connection).GetTCPConnection().Write([]byte("before ping .....\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error! ")
// 	}
// }

// // Test Handle
// func (this *PingRouter) Handle(request ziface.IRequest) {
// 	fmt.Println("Call PingRouter Handle ")
// 	//使用断言将接口类型转换为特定的结构体，然后访问结构体特有的元素，如果这里不用断言，用不了GetTCPConnection方法
// 	//GetConnection()方法返回的是接口类型，可以用断言转换为实现了该接口的结构体Connection
// 	//断言为*Connection类型，才可以调用
// 	_, err := request.GetConnection().(*znet.Connection).GetTCPConnection().Write([]byte("ping ...ping...ping\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error! ")
// 	}
// }

// // Test PostHandle
// func (this *PingRouter) PostHandle(request ziface.IRequest) {
// 	fmt.Println("Call Router PostHandle")
// 	_, err := request.GetConnection().(*znet.Connection).GetTCPConnection().Write([]byte("After ping .....\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error")
// 	}
// }

// // Server 模块的测试函数
// func main() {

// 	//1 创建一个server 句柄 s
// 	s := znet.NewServer()

// 	s.AddRouter(&PingRouter{})
// 	//2 开启服务
// 	s.Serve()
// }

// 只是负责测试datapack拆包，封包功能
func main() {
	//创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建服务器goroutine，负责从客户端goroutine读取粘包的数据，然后进行解析
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)
		}

		//处理客户端请求

		go func(conn net.Conn) {
			//创建封包拆包对象dp
			dp := znet.NewDataPack()

			for {
				//1、先读出流中的head部分
				headData := make([]byte, dp.GetHeadLen())

				//ReadFull  会把接收的空间（此处为headData）填充满为止，这个长度刚好读出包头
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error:", err) //错误为EOF（End Of File），意思为已经读到了文件末尾
					break
				}
				//2、将headData字节流，拆包到msg中
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpack err:", err)
					return
				}
				//3、根据msg.DataLen去继续在io流里面读取固定长度的字节作为Data
				if msgHead.GetDataLen() > 0 {
					//msg 是有data数据的，需要再次读取data数据
					msg := msgHead.(*znet.Message)
					//初始化Data空间
					msg.Data = make([]byte, msg.GetDataLen())

					//根据dataLen从io中读取字节流
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server un pack data err:", err)
						return
					}
					fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
				}

			}
		}(conn)
	}
}
