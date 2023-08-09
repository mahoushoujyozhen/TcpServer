package utils

import (
	"TcpServer/ziface"
	"encoding/json"
	"os"
)

/*
	存储一切有关框架的全局参数，供其他模块使用
	一些参数也可以通过 用户根据 zinx.json 来配置
*/

type GlobalObj struct {
	TcpServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听的端口
	Name      string         // 当前服务器名称
	Version   string         // 当前Zinx版本号

	MaxPacketSize uint32 //所需数据包的最大值
	MaxConn       int    // 当前服务主机允许最大的连接数
}

/*
	定义一个全局对象
*/

var GlobalObject *GlobalObj

//读取用户的配置文件

/*
在 Go 中，结构体方法可以使用两种类型的接收者：值接收者（p person）和指针接收者（p *person）
值接受者就等于值传递，一般用于get方法，不涉及实体值的修改
指针接受者就是指针传递，引用传递，涉及到实体值修改的时候使用
*/
func (g *GlobalObj) Reload() {
	//结构体方法中可以用g来访问调用这个方法的实体
	data, err := os.ReadFile("../conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析道struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供init方法，默认加载
*/

func init() {
	//初始化GlobalObj变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:          "TcpServer",
		Version:       "V0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}
	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
