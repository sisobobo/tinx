package tnet

import (
	"github.com/sisobobo/tinx/tlog"
	"math"
	"net"
)

const maxInt = math.MaxInt32

func initTCP(server *server, addrs []string, accept int) (err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			tlog.Errorf("net.ResolveTCPAddr(tcp, %s) error(%v)", bind, err)
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			tlog.Errorf("net.ListenTCP(tcp, %s) error(%v)", bind, err)
			return
		}
		tlog.Infof("start tcp listen: %s", bind)
		// split N core accept
		for i := 0; i < accept; i++ {
			go acceptTCP(server, listener)
		}
	}
	return
}

func acceptTCP(s *server, lis *net.TCPListener) {
	var (
		conn *net.TCPConn
		err  error
		r    int
	)
	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			// if listener close then return
			tlog.Errorf("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(s.conf.Server.KeepAlive); err != nil {
			tlog.Errorf("conn.SetKeepAlive() error(%v)", err)
			return
		}
		if err = conn.SetReadBuffer(s.conf.Server.RcvBuf); err != nil {
			tlog.Errorf("conn.SetReadBuffer() error(%v)", err)
			return
		}
		if err = conn.SetWriteBuffer(s.conf.Server.SndBuf); err != nil {
			tlog.Errorf("conn.SetWriteBuffer() error(%v)", err)
			return
		}
		ch := newChannel(s, conn, r).(*channel)
		go ch.open()
		if r++; r == maxInt {
			r = 0
		}
	}
}
