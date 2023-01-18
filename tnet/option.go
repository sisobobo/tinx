package tnet

type Option func(s *Server)

func SetCodec(codec Codec) Option {
	return func(s *Server) {
		s.codec = codec
	}
}

func SetHandler(handler Handler) Option {
	return func(s *Server) {
		s.handler = NewRouterHandler(handler)
	}
}
