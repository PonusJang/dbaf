package dbms

import (
	"dbaf/log"
	"dbaf/manager/models"
	"io"
	"net"
	"time"
)

func generateDBMS() (DBMS, func(io.Reader) ([]byte, error)) {
	return new(MySQL), MySQLReadPacket
}

func HandleClient(listenConn net.Conn, serverAddr *net.TCPAddr) error {
	d, reader := generateDBMS()
	log.Debug("客户端: %s", listenConn.RemoteAddr())
	serverConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Warn(err)
		listenConn.Close()
		return err
	}

	listenConn.SetReadDeadline(time.Now().Add(time.Minute * 1))
	serverConn.SetReadDeadline(time.Now().Add(time.Minute * 1))

	log.Debug("目标数据库: %s", serverConn.RemoteAddr())
	d.SetSockets(listenConn, serverConn)

	//设置证书
	//d.SetCertificate(config.Config.TLSCertificate, config.Config.TLSPrivateKey)
	d.SetReader(reader)
	err = d.Handler()
	if err != nil {
		log.Warn(err)
	}
	return err
}

type D struct {
	Id int
	F  *models.DbForward
}

type P struct {
	Id int
	S  *net.TCPAddr
	L  net.Listener
}

var (
	DChan      = make(chan D, 10)
	PChan      = make(chan P, 10)
	PCloseChan = make(chan P, 10)
)
