package znet

import "TcpServer/ziface"

// 实现router时,先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

//这里之所以BaseRouter的方法都为空，
//是因为有的Router不希望有PreHandle或者PostHandle，如果没有这个BaseRouter过度，那么之后的Router都需要去实现这两个用不着的方法
//所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHadle也可以实例子化

func (br *BaseRouter) PreHandle(req ziface.IRequest) {}

func (br *BaseRouter) handle(req ziface.IRequest) {}

func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
