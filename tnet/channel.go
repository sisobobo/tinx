package tnet

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sisobobo/tinx/tlog"
	"github.com/sisobobo/tinx/tpkg/bufio"
	"github.com/sisobobo/tinx/tpkg/bytes"
	"io"
	"net"
	"strings"
)

type Channel struct {
	id     string
	server *Server
	conn   net.Conn
	rp, wp *bytes.Pool
	rb, wb *bytes.Buffer
	reader bufio.Reader
	writer bufio.Writer
	bucket *Bucket
	r      int
}

func (c *Channel) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Channel) Id() string {
	return c.id

}

func (c *Channel) Close() {
	c.close()
}

func NewChannel(server *Server, conn *net.TCPConn, r int) *Channel {
	c := &Channel{
		id:     strings.ReplaceAll(uuid.NewString(), "-", ""),
		server: server,
		conn:   conn,
		rp:     server.round.Reader(r),
		wp:     server.round.Writer(r),
		bucket: server.round.Bucket(r),
	}
	c.rb = c.rp.Get()
	c.wb = c.wp.Get()
	c.reader.ResetBuffer(c.conn, c.rb.Bytes())
	c.writer.ResetBuffer(c.conn, c.wb.Bytes())
	c.bucket.Put(c)
	fmt.Printf("r:%d  %p : %d \n", r, c.bucket, c.bucket.ChannelCount())
	return c
}

func (c *Channel) open() {
	c.server.handler.Connect(c)
	go c.startReader()
}

func (c *Channel) startReader() {
	defer c.close()
	for {
		msg, err := c.server.codec.Decode(&c.reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			tlog.Errorf("%s接收到的据异常%s", c.RemoteAddr(), err)
			continue
		}
		go c.server.handler.Receive(c, msg)
	}
}

func (c *Channel) close() {
	c.bucket.Remove(c)
	c.rp.Put(c.rb)
	c.wp.Put(c.wb)
	c.conn.Close()
	c.server.handler.DisConnect(c)
}
