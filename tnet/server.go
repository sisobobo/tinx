package tnet

import (
	"github.com/sisobobo/tinx/tconf"
	"github.com/sisobobo/tinx/tlog"
	"math"
	"net"
)

type Server struct {
	conf      *tconf.Config
	serverID  string
	codec     Codec
	handler   Handler
	round     *Round
	bucketIdx uint32
}

func (s *Server) Start() {
	var (
		listener *net.TCPListener
		addr     *net.TCPAddr
		err      error
	)
	for _, bind := range s.conf.TCP.Bind {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			panic(err)
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			panic(err)
		}
		tlog.Infof("start tcp listen: %s", bind)
		for i := 0; i < 1; i++ {
			go s.acceptTcp(listener)
		}
	}
}

func (s *Server) acceptTcp(lis *net.TCPListener) {
	r := 0
	for {
		conn, err := lis.AcceptTCP()
		if err != nil {
			tlog.Errorf("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
			return
		}
		if err = conn.SetKeepAlive(s.conf.TCP.KeepAlive); err != nil {
			tlog.Errorf("conn.SetKeepAlive() error(%v)", err)
			return
		}
		if err = conn.SetReadBuffer(s.conf.TCP.RcvBuf); err != nil {
			tlog.Errorf("conn.SetReadBuffer() error(%v)", err)
			return
		}
		if err = conn.SetWriteBuffer(s.conf.TCP.SndBuf); err != nil {
			tlog.Errorf("conn.SetWriteBuffer() error(%v)", err)
			return
		}
		ch := NewChannel(s, conn, r)
		go ch.open()
		r++
		if r == math.MaxInt32 {
			r = 0
		}
	}
}

func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) Stop() {

}

func NewServer(configPath string, options ...Option) *Server {
	conf, err := tconf.NewConfig(configPath)
	if err != nil {
		panic(err)
	}
	s := &Server{
		conf:  conf,
		round: NewRound(conf),
	}
	s.setOptions(options...)
	return s
}

func (s *Server) setOptions(options ...Option) {
	for _, v := range options {
		v(s)
	}
}
