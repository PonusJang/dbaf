package dbms

import (
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
	"time"
)

const maxMSSQLPayloadLen = 4096

//MSSQL DBMS
type MSSQL struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   []byte
	username    []byte
	reader      func(io.Reader) ([]byte, error)
}

func (m *MSSQL) SetCertificate(crt, key string) (err error) {
	m.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
}

func (m *MSSQL) SetReader(f func(io.Reader) ([]byte, error)) {
	m.reader = f
}

func (m *MSSQL) SetSockets(c, s net.Conn) {
	m.client = c
	m.server = s
}

func (m *MSSQL) Close() {
	defer HandlePanic()
	m.client.Close()
	m.server.Close()
}

func (m *MSSQL) DefaultPort() uint {
	return 1433
}

func (m *MSSQL) Handler() error {
	defer HandlePanic()
	defer m.Close()

	success, err := m.handleLogin()
	if err != nil {
		return err
	}
	if !success {
		log.Warn("登录失败")
		return nil
	}
	for {
		var buf []byte
		buf, err = m.reader(m.client)
		if err != nil || len(buf) < 8 {
			return err
		}

		switch buf[0] {
		case 0x01: //SQL batch
			query := buf[8:]
			context := QueryContext{
				Query:    query,
				Database: m.currentDB,
				User:     m.username,
				Client:   RemoteAddrToIP(m.client.RemoteAddr()),
				Time:     time.Now(),
			}
			ProcessContext(context)
		}

		_, err = m.server.Write(buf)
		if err != nil {
			return err
		}

		err = ReadWrite(m.server, m.client, m.reader)
		if err != nil {
			return err
		}
	}
}

func (m *MSSQL) handleLogin() (success bool, err error) {

	buf, err := m.reader(m.client)
	if err != nil {
		return
	}

	if buf[0] != 0x12 {
		err = errors.New("该报文不是预登录请求报文")
		return
	}

	_, err = m.server.Write(buf)
	if err != nil {
		return
	}

	buf, err = m.reader(m.server)
	if err != nil {
		return
	}

	if buf[0] != 0x4 {
		err = errors.New("该报文不是预登录响应报文")
		return
	}

	_, err = m.client.Write(buf)
	if err != nil {
		return
	}

	data := buf[8:]

	var encryption byte

	for i := 0; i < len(data); i += 5 {
		switch data[i] {
		case 0x01: //Encryption
			encryption = data[int(data[i+1])]
			break
		case 0xff: //Terminator
			break
		}
	}
	log.Debug("Encryption: %v", encryption)

	buf, err = m.reader(m.client)
	if err != nil {
		return
	}

	m.username, m.currentDB = MSSQLGetUsernameDB(buf)

	for {
		//Receive PreLogin Request
		err = ReadWrite(m.client, m.server, m.reader)
		if err != nil {
			return
		}

		//Receive PRELOGIN response
		buf, err = m.reader(m.server)
		if err != nil {
			return
		}

		//Send PreLogin to server
		_, err = m.client.Write(buf)
		if err != nil {
			return
		}

		if buf[0] == 0x4 {
			break
		}
	}
	success = true
	return
}

var mssqlDataPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, maxMSSQLPayloadLen)
	},
}

func MSSQLReadPacket(src io.Reader) ([]byte, error) {
	data := mssqlDataPool.Get().([]byte)
	defer mssqlDataPool.Put(data)
	n, err := src.Read(data)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if n < maxMSSQLPayloadLen || err == io.EOF {
		return data[:n], nil
	}

	buf, err := MSSQLReadPacket(src)
	if err != nil {
		return nil, err
	}
	return append(data, buf...), nil
}

func MSSQLGetUsernameDB(data []byte) (username, db []byte) {
	//TODO: Extract Username and db name
	fmt.Printf("%x", data)
	return
}
