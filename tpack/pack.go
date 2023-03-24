package tpack

import (
	"encoding/binary"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
)

type BaseConnect struct {
}

func (d *BaseConnect) Connect(channel tiface.IChannel) {
	tlog.Infof("%s connect", channel.RemoteAddr())
}

func (d *BaseConnect) DisConnect(channel tiface.IChannel) {
	tlog.Infof("%s disconnect", channel.RemoteAddr())
}

type BaseMessage struct {
	msgId uint32
	data  []byte
}

func (d *BaseMessage) SetMsgId(id uint32) {
	d.msgId = id
}

func (d *BaseMessage) SetData(data []byte) {
	d.data = data
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

type DefaultPack struct {
	LengthFieldDecoder
}

func (d *DefaultPack) UnPack(channel tiface.IChannel, arr []byte) tiface.IMessage {
	msgId := d.order.Uint32(arr[:4])
	msg := NewBaseMessage(msgId, arr[4:])
	return msg
}

func (d *DefaultPack) Pack(channel tiface.IChannel, msg tiface.IMessage) []byte {
	data := msg.GetData()
	l := len(data) + 4
	lenArr := make([]byte, 2)
	msgIdArr := make([]byte, 4)
	d.order.PutUint16(lenArr, uint16(l))
	d.order.PutUint32(msgIdArr, msg.GetMsgId())
	header := append(lenArr, msgIdArr...)
	return append(header, data...)
}

func NewDefaultPack() tiface.IPack {
	return &DefaultPack{
		LengthFieldDecoder: NewLengthFieldDecoder(0, 2, 0, true, binary.BigEndian),
	}
}
