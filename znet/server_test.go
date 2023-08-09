package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// 模拟客户端
func ClientTest() {
	fmt.Println("Client Test ... start")
	//3秒后发器测试请求，给服务端开启服务的机会
	time.Sleep(time.Second * 3)
	//建立TCP连接
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}

	for {
		//往服务器发送数据
		_, err = conn.Write([]byte("hello ZINX"))
		if err != nil {
			fmt.Println("write error err", err)
			return
		}
		//读取服务器的响应
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf errpr")
			return
		}
		fmt.Printf("server call back :%s , cnt = %d \n", buf, cnt)
		time.Sleep(time.Second * 1)
	}
}

//测试入口函数，golang 会检测Test开头的函数，将其视为测试函数

func TestServer(t *testing.T) {

	/*
		服务端测试
	*/
	//1、创建一个server句柄
	s := NewServer()
	/*
		客户端测试
	*/

	go ClientTest()

	//2、开启服务
	s.Serve()
}
