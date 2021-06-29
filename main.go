package main

import (
	"dbaf/log"
	"dbaf/manager"
	"net"
)

func main() {

	serverAddr, _ := net.ResolveTCPAddr("tcp", "10.10.8.188:3307")
	l, err := net.Listen("tcp", "192.168.26.171:38888")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {

		listenConn, err := l.Accept()
		if err != nil {
			log.Warn("Error accepting connection: %v", err)
			continue
		}
		go handleClient(listenConn, serverAddr)
		go manager.RunServer()
	}

}
