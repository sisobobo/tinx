package tnet

import (
	"github.com/sisobobo/tinx/tconf"
	"github.com/sisobobo/tinx/tiface"
)

type server struct {
	conf          *tconf.Config
	round         *round
	buckets       []*bucket
	serverId      string
	connect       tiface.IConnect
	routerManager *routerManager
}

func (s *server) SetConnect(connect tiface.IConnect) {
	s.connect = connect
}

func (s *server) AddRouter(msgId uint32, router tiface.IRouter) {
	s.routerManager.addRouter(msgId, router)
}

func (s *server) Start() {
	if s.connect == nil {
		s.connect = newDefaultConnect()
	}
	if len(s.routerManager.routers) == 0 {
		panic("routers is nil , please add router")
	}
	err := initTCP(s, s.conf.Server.Bind, 1)
	if err != nil {
		panic(err)
	}
	select {}
}

func (s *server) Stop() {
	//TODO implement me
	panic("implement me")
}

func NewServer(c *tconf.Config) tiface.IServer {
	s := &server{
		conf:          c,
		round:         newRound(c),
		buckets:       make([]*bucket, c.Bucket.Size),
		routerManager: newRouterManager(),
	}
	for i := 0; i < c.Bucket.Size; i++ {
		s.buckets[i] = newBucket(c.Bucket)
	}
	//s.serverId = c.Env.Host
	return s
}
