package tnet

import (
	"context"
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
	id     string
	server *Server
	conn   net.Conn
	rp, wp *bytes.Pool
	rb, wb *bytes.Buffer
	reader bufio.Reader
	writer bufio.Writer
	bucket *Bucket
	ctx    context.Context
	cancel context.CancelFunc
	r      int
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

func (c *Channel) SendMessage(message Message) {
	c.RLock()
	defer c.RUnlock()
	data, err := c.server.codec.Encode(message)
	if err != nil {
		tlog.Errorf("encode error :", err)
		return
	}
	if _, err = c.writer.Write(data); err != nil {
		tlog.Errorf("write error :", err)
		return
	}
	if err = c.writer.Flush(); err != nil {
		tlog.Errorf("flush error :", err)
		return
	}
}

func (c *Channel) Close() {
	c.cancel()
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
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.server.handler.Connect(c)
	go c.startReader()
	select {
	case <-c.ctx.Done():
		c.close()
		return
	}
}

func (c *Channel) startReader() {
	defer c.Close()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msg, err := c.server.codec.Decode(&c.reader)
			if err == io.EOF {
				break
			}
			if err != nil {
				tlog.Errorf("decode error  :", err)
				return
			}
			go c.server.handler.Receive(c, msg)
		}
	}
}

func (c *Channel) close() {
	c.Lock()
	defer c.Unlock()
	c.bucket.Remove(c)
	c.rp.Put(c.rb)
	c.wp.Put(c.wb)
	c.conn.Close()
	c.server.handler.DisConnect(c)
}
