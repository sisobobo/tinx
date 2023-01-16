package tiface

import "context"

type Channel interface {
	LocalAddr() string
	RemoteAddr() string
	WriteAndFlush(interface{})
	Context() context.Context
	Close()
}
