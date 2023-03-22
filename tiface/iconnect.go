package tiface

import (
	"encoding/binary"
	"github.com/sisobobo/tinx/tpkg/bufio"
)

type IConnect interface {
	ByteOrder() binary.ByteOrder
	Connect(channel IChannel)
	Read(channel IChannel, reader *bufio.Reader) (IMessage, error)
	Write(channel IChannel, writer *bufio.Writer, message IMessage)
	DisConnect(channel IChannel)
}
