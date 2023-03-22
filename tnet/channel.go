package tnet

import (
	"context"
	"github.com/google/uuid"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
	"github.com/sisobobo/tinx/tpkg/bufio"
	"github.com/sisobobo/tinx/tpkg/bytes"
	"net"
	"strings"
	"sync"
)

type channel struct {
	sync.RWMutex
	id       string
	s        *server
	bucket   *bucket
	conn     *net.TCPConn
	rp, wp   *bytes.Pool
	rb, wb   *bytes.Buffer
	reader   bufio.Reader
	writer   bufio.Writer
	ctx      context.Context
	cancel   context.CancelFunc
	isClosed bool
}

func (ch *channel) Write(message tiface.IMessage) {
	ch.write(message)
}

func (ch *channel) Flush() {
	ch.writer.Flush()
}

func (ch *channel) WriteAndFlush(message tiface.IMessage) {
	ch.write(message)
	ch.writer.Flush()
}

func (ch *channel) Context() context.Context {
	return ch.ctx
}

func (ch *channel) RemoteAddr() net.Addr {
	return ch.conn.RemoteAddr()
}

func (ch *channel) LocalAddr() net.Addr {
	return ch.conn.LocalAddr()
}

func (ch *channel) Close() {
	ch.cancel()
}

func newChannel(s *server, conn *net.TCPConn, r int) tiface.IChannel {
	rp := s.round.reader(r)
	wp := s.round.writer(r)
	ch := &channel{
		id:     strings.ReplaceAll(uuid.NewString(), "-", ""),
		s:      s,
		conn:   conn,
		bucket: s.buckets[r],
		rp:     rp,
		wp:     wp,
		rb:     rp.Get(),
		wb:     wp.Get(),
	}
	ch.reader.ResetBuffer(ch.conn, ch.rb.Bytes())
	ch.writer.ResetBuffer(ch.conn, ch.wb.Bytes())
	ch.bucket.putChannel(ch)
	return ch
}

func (ch *channel) open() {
	ch.ctx, ch.cancel = context.WithCancel(context.Background())
	ch.s.connect.Connect(ch)
	go ch.startReader()
	select {
	case <-ch.ctx.Done():
		ch.close()
		return
	}
}

func (ch *channel) startReader() {
	defer ch.Close()
	for {
		select {
		case <-ch.ctx.Done():
			return
		default:
			msg, err := ch.s.connect.Read(ch, &ch.reader)
			if err != nil {
				return
			}
			if msg == nil {
				tlog.Warnf("[&s] msg is nil", ch.RemoteAddr())
				continue
			}
			ch.s.routerManager.route(ch, msg)
		}
	}
}

func (ch *channel) write(message tiface.IMessage) {
	if ch.isClosed {
		tlog.Errorf("send err , %s is closed", ch.RemoteAddr())
		return
	}
	ch.s.connect.Write(ch, &ch.writer, message)
}

func (ch *channel) close() {
	ch.s.connect.DisConnect(ch)
	ch.Lock()
	defer ch.Unlock()
	if ch.isClosed {
		return
	}
	ch.conn.Close()
	ch.bucket.delChannel(ch)
	ch.rp.Put(ch.rb)
	ch.wp.Put(ch.wb)
	ch.isClosed = true
}
