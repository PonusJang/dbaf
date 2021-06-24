package oracle

import (
	"bytes"
	"crypto/tls"
	flow "dbaf/analysis/io"
	"dbaf/analysis/sql"
	"dbaf/analysis/utils"
	"dbaf/log"
	"io"
	"net"
	"time"
)

type Oracle struct {
	client      net.Conn
	server      net.Conn
	certificate tls.Certificate
	currentDB   []byte
	username    []byte
	reader      func(io.Reader) ([]byte, error)
}

func (o *Oracle) SetCertificate(crt, key string) (err error) {
	o.certificate, err = tls.LoadX509KeyPair(crt, key)
	return
}

func (o *Oracle) SetReader(f func(io.Reader) ([]byte, error)) {
	o.reader = f
}

func (o *Oracle) SetSockets(c, s net.Conn) {
	defer utils.HandlePanic()
	o.client = c
	o.server = s
}

func (o *Oracle) Close() {
	defer utils.HandlePanic()
	o.client.Close()
	o.server.Close()
}

func (o *Oracle) DefaultPort() uint {
	return 1521
}

func (o *Oracle) Handler() error {
	defer utils.HandlePanic()
	defer o.Close()

	for {
		buf, err := o.readPacket(o.client)
		if err != nil {
			return err
		}
		var eof bool

		switch buf[4] { //Packet Type
		case 0x01: //Connect
			connectDataLen := int(buf[24])*256 + int(buf[25])
			connectData := buf[len(buf)-connectDataLen:]

			//Extracting Service name
			// FIXME: avoid string
			tmp1 := bytes.Split(connectData, []byte("SERVICE_NAME="))
			tmp2 := bytes.Split(tmp1[1], []byte{0x29}) // )
			o.currentDB = tmp2[0]

			log.Debug("Connect Data: %s", connectData)
			log.Debug("Service Name: %s", o.currentDB)
		case 0x06: //Data
			data := buf[8:]
			eof = data[1] == 0x40
			if !eof {
				payload := data[2:]
				if len(payload) > 16 && payload[0] == 0x11 && payload[15] == 0x03 && payload[16] == 0x5e {
					// I have no idea what this TTC is but its on top of query
					//simply skiping it
					payload = payload[15:]
				}
				switch payload[0] {
				case 0x03:
					switch payload[1] {
					case 0x5e: //reading query
						query, _ := sql.PascalString(payload[70:])
						context := sql.QueryContext{
							Query:    query,
							Database: o.currentDB,
							User:     o.username,
							Client:   utils.RemoteAddrToIP(o.client.RemoteAddr()),
							Time:     time.Now(),
						}
						sql.ProcessContext(context)
					case 0x76: // Reading username
						val, _ := sql.PascalString(payload[19:])
						o.username = val
						log.Debug("Username: %s", o.username)
					}
				}
			}
		}

		_, err = o.server.Write(buf)
		if err != nil || eof {
			return err
		}

		err = flow.ReadWrite(o.server, o.client, o.readPacket)
		if err != nil {
			return err
		}
	}
}

func (o *Oracle) readPacket(c io.Reader) (buf []byte, err error) {
	buf, err = o.reader(c)
	if err != nil {
		return
	}
	packetLen := int(buf[0])*256 + int(buf[1])
	for {
		if len(buf) == packetLen {
			break
		}
		var b []byte
		b, err = o.reader(c)
		if err != nil {
			return
		}
		buf = append(buf, b...)
	}
	return
}
