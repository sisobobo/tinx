package tnet

import "github.com/sisobobo/tinx/tiface"

type Message struct {
	handlerId interface{}
	msg       interface{}
}

func (m *Message) HandlerId() interface{} {
	return m.handlerId
}

func (m *Message) Msg() interface{} {
	return m.msg
}

func NewMessage(handlerKey interface{}, msg interface{}) tiface.IMessage {
	return &Message{
		handlerId: handlerKey,
		msg:       msg,
	}
}
