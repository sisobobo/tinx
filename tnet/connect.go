package tnet

import (
	"encoding/binary"
	"github.com/sisobobo/tinx/tiface"
	"github.com/sisobobo/tinx/tlog"
	"github.com/sisobobo/tinx/tpkg/bufio"
)

type BaseConnect struct {
}

func (c *BaseConnect) Connect(channel tiface.IChannel) {
	tlog.Infof("%s is connect", channel.RemoteAddr())
}

func (c *BaseConnect) DisConnect(channel tiface.IChannel) {
	tlog.Infof("%s is disconnect", channel.RemoteAddr())
}

type defaultConnect struct {
	BaseConnect
}

func (d *defaultConnect) ByteOrder() binary.ByteOrder {
	return binary.BigEndian
}

func (d *defaultConnect) Read(channel tiface.IChannel, reader *bufio.Reader) (tiface.IMessage, error) {
	header, err := reader.Pop(6)
	if err != nil {
		tlog.Errorf("read err : %s", err)
		return nil, err
	}
	l := d.ByteOrder().Uint16(header[:2]) & 8192
	msgId := d.ByteOrder().Uint32(header[2:])
	data, err := reader.Pop(int(l))
	if err != nil {
		tlog.Errorf("read err : %s", err)
		return nil, err
	}
	msg := NewBaseMessage(msgId, data)
	return msg, nil
}

func (d *defaultConnect) Write(channel tiface.IChannel, writer *bufio.Writer, message tiface.IMessage) {
	data := message.GetData()
	l := len(data) + 4
	header := make([]byte, 6)
	d.ByteOrder().PutUint16(header, uint16(l))
	d.ByteOrder().PutUint32(header, message.GetMsgId())
	bytes := append(header, data...)
	writer.Write(bytes)
}

func newDefaultConnect() tiface.IConnect {
	return &defaultConnect{}
}
