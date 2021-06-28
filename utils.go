package main

import (
	"dbaf/analysis/dbms"
	"dbaf/log"
	"io"
	"net"
)

func generateDBMS() (dbms.DBMS, func(io.Reader) ([]byte, error)) {
	return new(dbms.MySQL), dbms.MySQLReadPacket
}

func handleClient(listenConn net.Conn, serverAddr *net.TCPAddr) error {
	d, reader := generateDBMS()
	log.Debug("客户端: %s", listenConn.RemoteAddr())
	serverConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		log.Warn(err)
		listenConn.Close()
		return err
	}
	//if config.Config.Timeout > 0 {
	//	if err = listenConn.SetDeadline(time.Now().Add(config.Config.Timeout)); err != nil {
	//		return err
	//	}
	//	if err = serverConn.SetDeadline(time.Now().Add(config.Config.Timeout)); err != nil {
	//		return err
	//	}
	//}
	log.Debug("目标数据库: %s", serverConn.RemoteAddr())
	d.SetSockets(listenConn, serverConn)
	//d.SetCertificate(config.Config.TLSCertificate, config.Config.TLSPrivateKey)
	d.SetReader(reader)
	err = d.Handler()
	if err != nil {
		log.Warn(err)
	}
	return err
}
