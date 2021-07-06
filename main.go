package main

import (
	_ "dbaf/common"
	"dbaf/dbms"
	logger "dbaf/log"
	"dbaf/manager"
	db "dbaf/manager/databases"
	"dbaf/manager/models"
	"fmt"
	"net"
	"strconv"
)

func getServerAndListener(d *models.DbForward) (*net.TCPAddr, net.Listener) {
	serverAddr, _ := net.ResolveTCPAddr("tcp", d.DbIp+":"+strconv.Itoa(d.DbPort))
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(d.ListenPort))
	if err != nil {
		panic(err)
	}
	//defer l.Close()
	return serverAddr, l
}

func runInit() {
	var tmp []models.DbForward
	db.Db.Model(&models.DbForward{}).Find(&tmp)
	for i := 0; i < len(tmp); i++ {
		var d dbms.D
		d.Id = tmp[i].Id
		d.F = &tmp[i]
		dbms.DChan <- d
	}
}

func runProxy(l net.Listener, s *net.TCPAddr) {
	for {
		listenConn, err := l.Accept()
		if err != nil {
			logger.Warn("Error accepting connection: %v", err)
			continue
		}
		go dbms.HandleClient(listenConn, s)
	}

}

func main() {

	go manager.RunServer()

	go runInit()

	for {
		select {
		case d, ok := <-dbms.DChan:
			if ok {
				var p dbms.P
				p.S, p.L = getServerAndListener(d.F)
				p.Id = d.Id
				dbms.PChan <- p
				fmt.Print(d.F.ListenPort)
			}
			break
		case p, ok := <-dbms.PChan:
			if ok {
				go runProxy(p.L, p.S)
			}
			break
		case p, ok := <-dbms.PCloseChan:
			if ok {
				p.L.Close()
			}
			break
		default:
			//fmt.Printf("running.....")
		}

	}

}
