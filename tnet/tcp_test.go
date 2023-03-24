package tnet

import (
	"fmt"
	"github.com/sisobobo/tinx/tconf"
	"github.com/sisobobo/tinx/tiface"
	"testing"
)

type TRouter struct {
	BaseRouter
}

func (T TRouter) Handle(channel tiface.IChannel, msg tiface.IMessage) {
	fmt.Println(string(msg.GetData()))
	msg.SetMsgId(2)
	err := channel.Write(msg)
	if err != nil {
		fmt.Println("write err :", err)
	}
}

func TestTcp(t *testing.T) {
	conf := &tconf.Config{
		Bucket: &tconf.Bucket{
			Size:    32,
			Channel: 1024,
		},
		Server: &tconf.Server{
			Bind:         []string{":9999"},
			IsWs:         true,
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
