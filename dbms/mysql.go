package dbms

import (
	"bytes"
	"crypto/tls"
	"dbaf/common"
	"dbaf/log"
	"io"
	"net"
	"strings"
	"time"
)

const maxMySQLPayloadLen = 1<<24 - 1

type MySQL struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   []byte
	username    []byte
	reader      func(io.Reader) ([]byte, error)
}

func (m *MySQL) SetCertificate(crt, key string) (err error) {
	m.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
}

func (m *MySQL) SetReader(f func(io.Reader) ([]byte, error)) {
	m.reader = f
}

func (m *MySQL) SetSockets(c, s net.Conn) {
	m.client = c
	m.server = s
}

func (m *MySQL) Close() {
	defer common.HandlePanic()
	m.client.Close()
	m.server.Close()
}

func (m *MySQL) DefaultPort() uint {
	return 3306
}

func (m *MySQL) Handler() error {
	defer common.HandlePanic()
	defer m.Close()
	success, err := m.handleLogin()
	if err != nil {
		log.Error("处理登录失败")
		return err
	}
	if !success {
		log.Warn("登录失败")
		return nil
	}
	for {
		var buf []byte
		buf, err = ReadPacket(m.client)
		if err != nil || len(buf) < 5 {
			return err
		}
		data := buf[4:]
		switch data[0] {
		case 0x01: // 退出
			return nil
		case 0x02: //  use database
			m.currentDB = data[1:]
			log.Debug("使用数据库: %v", string(m.currentDB))
		case 0x03: //  查询
			query := data[1:]
			context := common.QueryContext{
				Query:    query,
				Database: m.currentDB,
				User:     m.username,
				Client:   common.RemoteAddrToIP(m.client.RemoteAddr()),
				Server:   common.RemoteAddrToIP(m.server.LocalAddr()),
				Time:     time.Now(),
			}
			common.ProcessContext(context)
		}

		// 转发
		_, err = m.server.Write(buf)

		if err != nil {
			log.Error("发送后端失败")
			return err
		}

		err = ReadWrite(m.server, m.client, m.reader)
		if err != nil {
			log.Error("转发失败")
			if err == io.EOF {
				continue
			}
			if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host.") {
				return err
			}
			log.Error(err.Error())
			return err
		}
	}
}

// 处理登录

func (m *MySQL) handleLogin() (success bool, err error) {

	err = ReadWrite(m.server, m.client, ReadPacket)
	if err != nil {
		log.Error("handleLogin失败")
		return
	}

	buf, err := ReadPacket(m.client)
	if err != nil {
		log.Error("handleLogin ReadPacket失败")
		return
	}
	data := buf[4:]

	m.username, m.currentDB = MySQLGetUsernameDB(data)

	//判断是否SSL
	ssl := (data[1] & 0x08) == 0x08

	//转发登录
	_, err = m.server.Write(buf)
	if err != nil {
		log.Error("handleLogin 转发登录失败")
		return
	}
	if ssl {
		m.client, m.server, err = common.TurnSSL(m.client, m.server, m.certificate)
		if err != nil {
			return
		}
		buf, err = ReadPacket(m.client)
		if err != nil {
			return
		}
		data = buf[4:]
		m.username, m.currentDB = MySQLGetUsernameDB(data)

		//Send Login Request
		_, err = m.server.Write(buf)
		if err != nil {
			log.Error("handleLogin 转发登录失败2")
			return
		}
	}
	log.Debug("使用SSL: %v", ssl)

	if len(m.currentDB) != 0 { //db Selected
		buf, err = ReadPacket(m.server)
		if err != nil {
			return
		}
	} else {
		//Receive Auth Switch Request
		err = ReadWrite(m.server, m.client, ReadPacket)
		if err != nil {
			return
		}
		//Receive Auth Switch Response
		err = ReadWrite(m.client, m.server, ReadPacket)
		if err != nil {
			return
		}
		//Receive Response Status
		buf, err = ReadPacket(m.server)
		if err != nil {
			return
		}
	}

	if buf[5] != 0x15 {
		success = true
	}

	// 往客户端发送响应
	_, err = m.client.Write(buf)
	return
}

func MySQLGetUsernameDB(data []byte) (username, db []byte) {
	if len(data) < 33 {
		return
	}
	pos := 32

	nullByteIndex := bytes.IndexByte(data[pos:], 0x00)
	username = data[pos : nullByteIndex+pos]
	log.Debug("用户名: %s", string(username))
	pos += nullByteIndex + 22
	nullByteIndex = bytes.IndexByte(data[pos:], 0x00)

	dbSelectedCheck := len(data) > nullByteIndex+pos+1

	if nullByteIndex != 0 && dbSelectedCheck {
		db = data[pos : nullByteIndex+pos]
		log.Debug("数据库: %s", string(db))
	}
	return
}

func MySQLReadPacket(src io.Reader) ([]byte, error) {
	data := make([]byte, maxMySQLPayloadLen)
	_, err := src.Read(data)
	return data, err
}
