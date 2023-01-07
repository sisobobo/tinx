package tnet

import (
	"bufio"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
	"io"
	"net"
)

type Channel struct {
	id     uint32        //Id
	server *Server       //server
	conn   *net.TCPConn  //conn
	reader *bufio.Reader // reader
	writer *bufio.Writer // writer
}

func (c *Channel) WriteAndFlush(message tiface.Message) {
	data, err := c.server.pack.Encode(message)
	if err != nil {
		tlog.Error("encode error : %s ", err)
		return
	}
	c.writer.Write(data)
	c.server.pack.Pack(c.writer)
	c.writer.Flush()
}

func (c *Channel) Id() uint32 {
	return c.id
}

func (c *Channel) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Channel) LocalAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func NewChannel(server *Server, conn *net.TCPConn) tiface.IChannel {
	maxFrameLength := int(server.pack.GetMaxFrameLength())
	c := &Channel{
		server: server,
		conn:   conn,
		reader: bufio.NewReaderSize(conn, maxFrameLength),
		writer: bufio.NewWriterSize(conn, maxFrameLength),
	}
	return c
}

func (c *Channel) open() {
	tlog.INFO("连接：%s", c.conn.RemoteAddr().String())
	go c.startReader()
	go c.startWriter()
}

func (c *Channel) close() {
	tlog.INFO("连接断开")
	c.conn.Close()
}

func (c *Channel) startReader() {
	defer c.close()
	for {
		data, err := c.server.pack.UnPack(c.reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		if len(data) > 0 {
			key, msg, err := c.server.pack.Decode(data)
			if err != nil {
				continue
			}
			if msg != nil {
				go c.server.handlerManager.doMsgHandler(c, key, msg)
			}
		}
	}
}

func (c *Channel) startWriter() {
	//c.server.pack.Encode()
}
