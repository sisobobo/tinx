package tnet

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
	"github.com/sisobobo/tinx/tpkg/bufio"
	"github.com/sisobobo/tinx/tpkg/bytes"
	"github.com/sisobobo/tinx/tpkg/websocket"
	"io"
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
	ws       *websocket.Conn
	rp, wp   *bytes.Pool
	rb, wb   *bytes.Buffer
	reader   bufio.Reader
	writer   bufio.Writer
	ctx      context.Context
	cancel   context.CancelFunc
	isClosed bool
}

func (ch *channel) Flush() error {
	return ch.writer.Flush()
}

func (ch *channel) WriteAndFlush(message tiface.IMessage) error {
	err := ch.write(message)
	if err != nil {
		return err
	}
	return ch.Flush()
}

func (ch *channel) Write(message tiface.IMessage) error {
	return ch.write(message)
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
	ch.s.pack.Connect(ch)
	if ch.s.conf.Server.IsWs {
		go ch.startWsReader()
	} else {
		go ch.startTcpReader()
	}
	select {
	case <-ch.ctx.Done():
		ch.close()
		return
	}
}

func (ch *channel) startWsReader() {
	defer ch.Close()
	req, err := websocket.ReadRequest(&ch.reader)
	if err != nil || req.RequestURI != "/ws" {
		if err != io.EOF {
			tlog.Errorf("http.ReadRequest(rr) error(%v)", err)
		}
		return
	}
	ws, err := websocket.Upgrade(ch.conn, &ch.reader, &ch.writer, req)
	if err != nil {
		if err != io.EOF {
			tlog.Errorf("http.ReadRequest(rr) error(%v)", err)
		}
		return
	}
	ch.ws = ws
	for {
		select {
		case <-ch.ctx.Done():
			return
		default:
			_, data, err := ws.ReadMessage()
			if err != nil {
				return
			}
			reader := bufio.NewReaderSize(bytes.NewReader(data), len(data))
			err = ch.router(reader)
			if err != nil {
				return
			}
		}
	}
}

func (ch *channel) router(reader *bufio.Reader) error {
	arr, err := ch.s.pack.Read(ch, reader)
	if err != nil {
		return err
	}
	msg := ch.s.pack.UnPack(ch, arr)
	if msg == nil {
		tlog.Warnf("[%s] msg is nil", ch.RemoteAddr())
		return nil
	}
	ch.s.routerManager.route(ch, msg)
	return nil
}

func (ch *channel) startTcpReader() {
	defer ch.Close()
	for {
		select {
		case <-ch.ctx.Done():
			return
		default:
			err := ch.router(&ch.reader)
			if err != nil {
				return
			}
		}
	}
}

func (ch *channel) write(message tiface.IMessage) error {
	var err error
	if ch.isClosed {
		return errors.New("connect is closed")
	}
	arr := ch.s.pack.Pack(ch, message)
	if ch.s.conf.Server.IsWs {
		err = ch.ws.WriteMessage(websocket.BinaryMessage, arr)
	} else {
		_, err = ch.writer.Write(arr)
	}
	return err
}

func (ch *channel) close() {
	ch.s.pack.DisConnect(ch)
	ch.Lock()
	defer ch.Unlock()
	if ch.isClosed {
		return
	}
	if ch.ws != nil {
		ch.ws.Close()
	} else {
		ch.conn.Close()
	}
	ch.bucket.delChannel(ch)
	ch.rp.Put(ch.rb)
	ch.wp.Put(ch.wb)
	ch.isClosed = true
}
