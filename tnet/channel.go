package tnet

import (
	"context"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tpkg/bufio"
	"github.com/sisobobo/tinx/tpkg/bytes"
	"io"
	"net"
)

type Channel struct {
	id     uint32
	server *Server
	conn   net.Conn
	writer bufio.Writer
	reader bufio.Reader
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *Channel) WriteAndFlush(msg interface{}) {
	data, err := c.server.codec.Encode(msg)
	if err != nil {
		return
	}
	c.writer.Write(data)
	c.writer.Flush()
}

func (c *Channel) Close() {
	c.cancel()
}

func (c *Channel) Context() context.Context {
	return c.ctx
}

func (c *Channel) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

func (c *Channel) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func newChannel(server *Server, conn net.Conn, rb, wb *bytes.Buffer) tiface.Channel {
	c := &Channel{
		server: server,
		conn:   conn,
	}
	c.writer.ResetBuffer(conn, wb.Bytes())
	c.reader.ResetBuffer(conn, rb.Bytes())
	return c
}

func (c *Channel) open() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.server.handler.Active(c)
	go c.read()
	select {
	case <-c.ctx.Done():
		c.close()
		return
	}
}

func (c *Channel) read() {
	defer c.close()
	for {
		data, err := c.server.codec.Decode(&c.reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		c.server.handler.Read(c, data)
	}
}

func (c *Channel) close() {
	c.server.handler.InActive(c)
	c.conn.Close()
}
