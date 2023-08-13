package main

import (
	"TcpServer/znet"
	"fmt"
	"net"
)

// func main() {

// 	fmt.Println("Client Test ... start")
// 	//3秒之后发起测试请求，给服务端开启服务的机会
// 	time.Sleep(3 * time.Second)

// 	conn, err := net.Dial("tcp", "127.0.0.1:7777")
// 	if err != nil {
// 		fmt.Println("client start err, exit!")
// 		return
// 	}

// 	for {
// 		_, err := conn.Write([]byte("hahaha"))
// 		if err != nil {
// 			fmt.Println("write error err ", err)
// 			return
// 		}

// 		buf := make([]byte, 512)
// 		cnt, err := conn.Read(buf)
// 		if err != nil {
// 			fmt.Println("read buf error ")
// 			return
// 		}

// 		fmt.Printf(" server call back : %s, cnt = %d\n", buf, cnt)

// 		time.Sleep(1 * time.Second)
// 	}
// }

func main() {
	//  客户端goroutine，负责模拟粘包的数据，然后进行发送
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	//创建一个封包对象dp
	dp := znet.NewDataPack()

	//封装一个msg1包
	msg1 := &znet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err :", err)
		return
	}

	msg2 := &znet.Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}

	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}

	//将sendData1 和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	//客户端阻塞
	select {}

}
