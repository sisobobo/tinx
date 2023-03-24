package tiface

import (
	"context"
	"net"
)

type IChannel interface {
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	Context() context.Context
	Write(message IMessage) error
	Flush() error
	WriteAndFlush(message IMessage) error
	Close()
}
