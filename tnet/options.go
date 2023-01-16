package tnet

import "github.com/sisobobo/tinx/tiface"

type Option func(s *Server)

func SetHandler(handler tiface.Handler) Option {
	return func(s *Server) {
		s.handler = handler
	}
}

func SetCodec(codec tiface.Codec) Option {
	return func(s *Server) {
		s.codec = codec
	}
}
