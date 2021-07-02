package main

import (
	_ "dbaf/common"
	"dbaf/dbms"
	"dbaf/manager"
	"dbaf/manager/models"
	"fmt"
	"net"
	"strconv"
)

func getServerAndListener(d *models.DbForward) (*net.TCPAddr, net.Listener) {
	serverAddr, _ := net.ResolveTCPAddr("tcp", d.DbIp+":"+strconv.Itoa(d.DbPort))
	l, err := net.Listen("tcp", "192.168.26.171:"+strconv.Itoa(d.ListenPort))
	if err != nil {
		panic(err)
	}
	//defer l.Close()
	return serverAddr, l
}

func main() {

	go manager.RunServer()

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
				l, err := p.L.Accept()
				if err != nil {
					dbms.PCloseChan <- p
				}
				go dbms.HandleClient(l, p.S)
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
