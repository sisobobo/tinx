package tnet

import (
	"github.com/sisobobo/tinx/tconf"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
	"github.com/sisobobo/tinx/tpack"
	"math"
	"net"
)

type server struct {
	conf          *tconf.Config
	round         *round
	buckets       []*bucket
	serverId      string
	pack          tiface.IPack
	routerManager *routerManager
	exitChan      chan bool
}

func (s *server) SetPack(pack tiface.IPack) {
	s.pack = pack
}

func (s *server) AddRouter(msgId uint32, router tiface.IRouter) {
	s.routerManager.addRouter(msgId, router)
}

func (s *server) Start() {
	if s.pack == nil {
		s.pack = tpack.NewDefaultPack()
	}
	if len(s.routerManager.routers) == 0 {
		panic("routers is nil , please add router")
	}
	lis, err := s.initServer(s.conf.Server.Bind, 1)
	if err != nil {
		panic(err)
	}
	select {
	case <-s.exitChan:
		for _, v := range lis {
			err := v.Close()
			if err != nil {
				tlog.Errorf("Listener close err :%v", err)
			}
			tlog.Infof("listener %s closed", v.Addr().String())
		}
	}
}

func (s *server) Stop() {
	tlog.Infof("server closed")
	s.exitChan <- true
	close(s.exitChan)
}

func NewServer(c *tconf.Config) tiface.IServer {
	s := &server{
		conf:          c,
		round:         newRound(c),
		buckets:       make([]*bucket, c.Bucket.Size),
		routerManager: newRouterManager(),
		exitChan:      make(chan bool),
	}
	for i := 0; i < c.Bucket.Size; i++ {
		s.buckets[i] = newBucket(c.Bucket)
	}
	//s.serverId = c.Env.Host
	return s
}

func (s *server) initServer(addrs []string, accept int) (lis []*net.TCPListener, err error) {
	var (
		bind     string
		listener *net.TCPListener
		addr     *net.TCPAddr
	)
	lis = make([]*net.TCPListener, len(addrs))
	for i := 0; i < len(addrs); i++ {
		bind = addrs[i]
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
			go acceptTCP(s, listener)
		}
		lis[i] = listener
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
		if r++; r == math.MaxInt32 {
			r = 0
		}
	}
}
