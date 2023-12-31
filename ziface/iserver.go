package ziface

// 定义服务器接口
type IServer interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	// 开启业务服务方法
	Serve()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链路使用
	AddRouter(msgId uint32, router IRouter)

	//得到连接管理器
	GetConnMgr() IConnManager

	//设置该Server的链接创建时Hook函数
	SetOnConnStart(func(IConnection))
	//设置该Server链接断开时的Hook函数
	SetOnConnStop(func(IConnection))
	//调用链接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//嗲用链接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)
}
