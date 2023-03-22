package tnet

import (
	"github.com/sisobobo/tinx/tconf"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tpkg/bufio"
	"testing"
)

type MyConnect struct {
	BaseConnect
}

func (m *MyConnect) Read(channel tiface.IChannel, reader *bufio.Reader) ([]byte, error) {
	pop, err := reader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	return pop[:len(pop)-2], nil
}

type TRouter struct {
	BaseRouter
}

func (T TRouter) Handle(channel tiface.IChannel, msg tiface.IMessage) {
}

func TestTcp(t *testing.T) {
	conf := &tconf.Config{
		Bucket: &tconf.Bucket{
			Size:    32,
			Channel: 1024,
		},
		Server: &tconf.Server{
			Bind:         []string{":3101"},
			SndBuf:       4096,
			RcvBuf:       4096,
			KeepAlive:    false,
			Reader:       2,
			ReadBuf:      1024,
			ReadBufSize:  8192,
			Writer:       2,
			WriteBuf:     1024,
			WriteBufSize: 8192,
		},
	}
	s := NewServer(conf)
	s.AddRouter(1, &TRouter{})
	s.Start()
}
