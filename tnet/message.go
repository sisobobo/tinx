package tnet

import "github.com/sisobobo/tinx/tiface"

type BaseMessage struct {
	msgId uint32
	data  []byte
}

func (d *BaseMessage) GetMsgId() uint32 {
	return d.msgId
}

func (d *BaseMessage) GetData() []byte {
	return d.data
}

func NewBaseMessage(msgId uint32, data []byte) tiface.IMessage {
	return &BaseMessage{
		msgId: msgId,
		data:  data,
	}
}
