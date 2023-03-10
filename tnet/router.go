package tnet

import (
	"fmt"
	"github.com/sisobobo/tinx/tlog"
)

type RouterId interface{}

type Router interface {
	ID() RouterId
	Handler(channel *Channel, msg Message)
}

type RouterManager struct {
	routers map[interface{}]Router
}

func NewRouterManager() *RouterManager {
	return &RouterManager{
		routers: make(map[interface{}]Router),
	}
}

func (rm *RouterManager) add(id interface{}, r Router) {
	router := rm.routers[id]
	if router != nil {
		panic(fmt.Sprintf("routerId %v is already a route", id))
	}
	rm.routers[id] = r
}

func (rm *RouterManager) doHandler(channel *Channel, msg Message) {
	if msg.RouterId() != nil {
		router := rm.routers[msg.RouterId()]
		if router == nil {
			tlog.Warnf("routerId %v not find router", msg.RouterId())
			return
		}
		go router.Handler(channel, msg)
	}
}
