package tnet

import (
	"fmt"
	"net"
	"tinx/tiface"
	"tinx/tlog"
)

type Server struct {
	Name           string
	Host           string
	Port           int
	pack           tiface.IPack
	handlerManager *HandlerManager
}

func (s *Server) AddMsgHandlers(handlers ...tiface.IHandler) {
	for _, handler := range handlers {
		s.handlerManager.addHandler(handler.Id(), handler)
	}
}

func (s *Server) preStart() {
	if s.pack == nil {
		panic("pack is nil , please check SetPack()")
	}
}

func (s *Server) initTcp() (listener *net.TCPListener, err error) {
	var (
		address = fmt.Sprintf("%s:%d", s.Host, s.Port)
		addr    *net.TCPAddr
	)
	if addr, err = net.ResolveTCPAddr("tcp", address); err != nil {
		return
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return
	}
	tlog.INFO("start up listen:%s", address)
	return
}

func (s *Server) Start() {
	s.preStart()
	lis, err := s.initTcp()
	if err != nil {
		tlog.Error("server start error : %s", err)
		return
	}
	go s.acceptTcp(lis)
	select {}
}

func (s *Server) acceptTcp(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			tlog.Error("listener.Accept(\"%s\") error(%v)", listener.Addr().String(), err)
			return
		}
		channel := NewChannel(s, conn).(*Channel)
		channel.open()
	}
}

func (s *Server) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) SetPack(pack tiface.IPack) {
	s.pack = pack
}

func NewServer(name string, port int) tiface.IServer {
	return &Server{
		Name:           name,
		Host:           "0.0.0.0",
		Port:           port,
		handlerManager: newHandlerManager(),
	}
}
