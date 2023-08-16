package znet

import (
	"TcpServer/ziface"
)

type Request struct {
	//conn这个句柄建立了 request和connection的相互索引的关系
	conn ziface.IConnection // 已经和客户端建立好的连接  句柄，标志该request属于哪个connection，在request可以使用connection
	msg  ziface.IMessage    //客户端请求的数据
}

//获取请求连接信息

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求消息的ID
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
