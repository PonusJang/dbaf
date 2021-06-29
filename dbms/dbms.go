package dbms

import (
	"io"
	"net"
)

// 数据库接口

type DBMS interface {
	DefaultPort() uint                         // 数据库端口
	Close()                                    //关闭
	SetReader(func(io.Reader) ([]byte, error)) //读取数据
	Handler() error                            //处理函数
	SetSockets(net.Conn, net.Conn)             //客户端 后端数据库
	SetCertificate(string, string) error       //设置SSL证书
}
