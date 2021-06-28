package dbms

import (
	"bytes"
	"crypto/tls"
	"dbaf/log"
	"errors"
	"io"
	"net"
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
	defer HandlePanic()
	m.client.Close()
	m.server.Close()
}

func (m *MySQL) DefaultPort() uint {
	return 3306
}

func (m *MySQL) Handler() error {
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
			log.Debug("使用数据库: %v", m.currentDB)
		case 0x03: //  查询
			query := data[1:]
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

// 处理登录

func (m *MySQL) handleLogin() (success bool, err error) {

	err = ReadWrite(m.server, m.client, ReadPacket)
	if err != nil {
		return
	}

	buf, err := ReadPacket(m.client)
	if err != nil {
		return
	}
	data := buf[4:]

	m.username, m.currentDB = MySQLGetUsernameDB(data)

	//判断是否SSL
	ssl := (data[1] & 0x08) == 0x08

	//转发登录
	_, err = m.server.Write(buf)
	if err != nil {
		return
	}
	if ssl {
		m.client, m.server, err = TurnSSL(m.client, m.server, m.certificate)
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
			return
		}
	}
	log.Debug("SSL bit: %v", ssl)

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

// 读取报文

func MySQLReadPacket(src io.Reader) ([]byte, error) {
	data := make([]byte, maxMySQLPayloadLen)
	var prevData []byte
	for {

		n, err := src.Read(data)
		if err != nil {
			return nil, err
		}
		data = data[:n]
		pktLen := int(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16)

		if pktLen == 0 {
			if prevData == nil {
				return nil, errors.New("Malform Packet")
			}

			return prevData, nil
		}

		eof := true
		if len(data) > 8 {
			tail := data[len(data)-9:]
			eof = tail[0] == 5 && tail[1] == 0 && tail[2] == 0 && tail[4] == 0xfe
		}

		if eof {
			if prevData == nil {
				return data, nil
			}

			return append(prevData, data...), nil
		}

		prevData = append(prevData, data...)
	}
}

func MySQLGetUsernameDB(data []byte) (username, db []byte) {
	if len(data) < 33 {
		return
	}
	pos := 32

	nullByteIndex := bytes.IndexByte(data[pos:], 0x00)
	username = data[pos : nullByteIndex+pos]
	log.Debug("Username: %s", username)
	pos += nullByteIndex + 22
	nullByteIndex = bytes.IndexByte(data[pos:], 0x00)

	dbSelectedCheck := len(data) > nullByteIndex+pos+1

	if nullByteIndex != 0 && dbSelectedCheck {
		db = data[pos : nullByteIndex+pos]
		log.Debug("Database: %s", db)
	}
	return
}