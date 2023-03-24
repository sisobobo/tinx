package tnet

import (
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
	"strconv"
)

type routerManager struct {
	routers map[uint32]tiface.IRouter
}

func newRouterManager() *routerManager {
	return &routerManager{
		routers: make(map[uint32]tiface.IRouter),
	}
}

func (rm *routerManager) addRouter(msgId uint32, router tiface.IRouter) {
	if _, ok := rm.routers[msgId]; ok {
		panic("repeated router msgId:" + strconv.Itoa(int(msgId)))
	}
	rm.routers[msgId] = router
}

func (rm *routerManager) route(channel tiface.IChannel, msg tiface.IMessage) {
	router, ok := rm.routers[msg.GetMsgId()]
	if !ok {
		tlog.Warnf("router msgID = %d is not found", msg.GetMsgId())
		return
	}
	router.PreHandle(channel, msg)
	router.Handle(channel, msg)
	router.PostHandle(channel, msg)
}

type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(channel tiface.IChannel, msg tiface.IMessage) {

}

func (b *BaseRouter) Handle(channel tiface.IChannel, msg tiface.IMessage) {

}

func (b *BaseRouter) PostHandle(channel tiface.IChannel, msg tiface.IMessage) {
}
