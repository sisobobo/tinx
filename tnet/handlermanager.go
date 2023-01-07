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

func (m *HandlerManager) doMsgHandler(channel tiface.IChannel, key tiface.HandlerKey, message tiface.Message) {
	handler, ok := m.handlers[key]
	if !ok {
		tlog.WARN("handler %v is not found !", key)
		return
	}
	handler.PreHandler(channel, message)
	handler.Handler(channel, message)
	handler.PostHandler(channel, message)
}

func (m *HandlerManager) addHandler(key interface{}, handler tiface.IHandler) {
	//如果已经添加，返回
	if _, ok := m.handlers[key]; ok {
		tlog.Error("key %v had a handler", key)
		panic("add handler repeatedly")
	}
	m.handlers[key] = handler
}
