package tiface

import "net"

type IChannel interface {
	Id() uint32
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
}
