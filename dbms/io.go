package dbms

import (
	"bytes"
	"io"
	"net"
)

const (
	chunkSize = 4096
)

//读取报文

func ReadPacket(conn io.Reader) ([]byte, error) {
	data := make([]byte, chunkSize)
	buf := bytes.Buffer{}
	for {
		n, err := conn.Read(data)
		if err != nil {
			return nil, err
		}
		buf.Write(data[:n])
		if n != chunkSize {
			break
		}
	}
	return buf.Bytes(), nil
}

// 将请求转发后端数据库

func ReadWrite(src, dst net.Conn, reader func(io.Reader) ([]byte, error)) error {

	buf, err := reader(src)

	if err != nil {
		return err
	}

	_, err = dst.Write(buf)
	return err
}
