package tnet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sisobobo/tinx/tlog"
	"github.com/sisobobo/tinx/tpkg/bufio"
	"github.com/sisobobo/tinx/tpkg/bytes"
	"io"
	"net"
	"strings"
	"sync"
)

type Channel struct {
	sync.RWMutex
	id       string
	server   *Server
	conn     net.Conn
	rp, wp   *bytes.Pool
	rb, wb   *bytes.Buffer
	reader   bufio.Reader
	writer   bufio.Writer
	bucket   *Bucket
	isClosed bool
	r        int
}

func (c *Channel) Id() string {
	return c.id
}

func (c *Channel) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Channel) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

func (c *Channel) SendMessage(message Message) error {
	data, err := c.server.codec.Encode(message)
	if err != nil {
		return err
	}
	if c.isClosed {
		return errors.New("connection is closed")
	}
	if _, err = c.writer.Write(data); err != nil {
		return err
	}
	if err = c.writer.Flush(); err != nil {
		return err
	}
	return nil
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
	return c
}

func (c *Channel) open() {
	c.isClosed = false
	c.server.handler.Connect(c)
	defer c.Close()
	for {
		msg, err := c.server.codec.Decode(&c.reader)
		if c.isClosed {
			return
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			tlog.Errorf("decode error : %s", err)
			return
		}
		c.server.handler.Receive(c, msg)
	}
}

func (c *Channel) close() {
	c.Lock()
	defer c.Unlock()
	if c.isClosed {
		return
	}
	c.bucket.Remove(c)
	c.rp.Put(c.rb)
	c.wp.Put(c.wb)
	c.conn.Close()
	c.server.handler.DisConnect(c)
	c.isClosed = true
}
