package tnet

import (
	"fmt"
	"github.com/sisobobo/tinx/tiface"
	"math"
	"net"
)

type Server struct {
	id         string
	Ip         string
	port       uint16
	NetWorking string
	round      *Round
	codec      tiface.Codec
	handler    tiface.Handler
}

func (s *Server) Stop() {

}

func NewServer(id string, port uint16, options ...Option) tiface.Server {
	s := &Server{
		id:         id,
		port:       port,
		NetWorking: "tcp",
		round:      newRound(),
	}
	s.setOptions(options...)
	return s
}

func (s *Server) Start() {
	lis := s.listener()
	go connect(s, lis)
}

func (s *Server) Serve() {
	if s.handler == nil {
		panic("handler not allow nil")
	}
	if s.codec == nil {
		panic("codec not allow nil")
	}
	s.Start()
	select {}
}

func connect(s *Server, lis net.Listener) {
	n := 0
	for {
		conn, err := lis.Accept()
		if err != nil {
			conn.Close()
			return
		}
		rp := s.round.Reader(n)
		wp := s.round.Writer(n)
		ch := newChannel(s, conn, rp.Get(), wp.Get()).(*Channel)
		ch.open()
		if n++; n == math.MaxInt32 {
			n = 0
		}
	}
}

func (s *Server) listener() net.Listener {
	address := fmt.Sprintf("%s:%d", s.Ip, s.port)
	listen, err := net.Listen(s.NetWorking, address)
	if err != nil {
		panic(err)
	}
	return listen
}

func (s *Server) setOptions(options ...Option) {
	for i := 0; i < len(options); i++ {
		options[i](s)
	}
}
