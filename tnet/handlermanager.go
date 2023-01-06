package tnet

import (
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
)

type HandlerManager struct {
	handlers map[interface{}]tiface.IHandler
}

func newHandlerManager() *HandlerManager {
	return &HandlerManager{
		handlers: map[interface{}]tiface.IHandler{},
	}
}

func (m *HandlerManager) doMsgHandler(channel tiface.IChannel, message tiface.IMessage) {
	handler, ok := m.handlers[message.HandlerId()]
	if !ok {
		tlog.WARN("handler %v is not found !", message.HandlerId())
		return
	}
	handler.PreHandler(channel, message.Msg())
	handler.Handler(channel, message.Msg())
	handler.PostHandler(channel, message.Msg())
}

func (m *HandlerManager) addHandler(key interface{}, handler tiface.IHandler) {
	//如果已经添加，返回
	if _, ok := m.handlers[key]; ok {
		tlog.Error("key %v had a handler", key)
		panic("add handler repeatedly")
	}
	m.handlers[key] = handler
}
