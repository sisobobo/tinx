package tnet

import (
	"bufio"
	"io"
	"net"
	"tinx/tiface"
	"tinx/tlog"
)

type Channel struct {
	id     uint32        //Id
	server *Server       //server
	conn   *net.TCPConn  //conn
	reader *bufio.Reader // reader
	writer *bufio.Writer // writer
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
	c := &Channel{
		server: server,
		conn:   conn,
	}
	return c
}

func (c *Channel) open() {
	tlog.INFO("连接：%s", c.conn.RemoteAddr().String())
	go c.startReader()
}

func (c *Channel) close() {
	tlog.INFO("连接断开")
	c.conn.Close()
}

func (c *Channel) startReader() {
	c.reader = bufio.NewReaderSize(c.conn, int(c.server.pack.GetMaxFrameLength()))
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
			msg, err := c.server.pack.Decode(data)
			if err != nil {
				continue
			}
			if msg != nil {
				go c.server.handlerManager.doMsgHandler(c, msg)
			}
		}
	}
}
