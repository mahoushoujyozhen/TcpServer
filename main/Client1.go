package main

import (
	"TcpServer/znet"
	"fmt"
	"io"
	"net"
	"time"
)

/*
	模拟客户端
*/

func main() {
	fmt.Println("Client Test ... start")
	//3秒后发起测试请求，给服务端开启服务端机会
	time.Sleep(time.Second * 3)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}

	for {
		//发封包message消息
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx V0.6 Clent1 Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		/*
			读server返回的响应数据
		*/
		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error")
			break
		}

		//将headData字节流，拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err: ", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message) //接口类型断言成具体结构，才能够拿到结构体内的Data
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}
			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))

		}
		time.Sleep(time.Second * 1)
	}
}
