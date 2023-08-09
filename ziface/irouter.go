package ziface

/*
	路由接口，这里的路由是 使用框架或者给该链接自定的 处理业务方法
	路由里的IRequest 则包含用该链接的链接信息和该链接的请求数据信息
*/

type IRouter interface {
	PreHandle(requset IRequest)  //在处理conn业务之前的钩子方法
	handle(request IRequest)     //处理conn业务的方法
	PostHandle(request IRequest) //处理conn业务之后的钩子函数
}
