package dbms

import (
	"bytes"
	logger "dbaf/log"
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
	logger.Debug("读取......")
	logger.Debug(len(buf))
	if err != nil {
		logger.Debug("读取...出错了...")
		return err
	}
	logger.Debug("写入...........")
	_, err = dst.Write(buf)
	if err != nil {
		logger.Debug("写入...出错了...")
		return err
	}
	return err

}
