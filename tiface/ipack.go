package tiface

import (
	"bufio"
)

type HandlerKey any
type Message any

type IPack interface {
	GetMaxFrameLength() uint32
	Pack(writer *bufio.Writer) error
	UnPack(reader *bufio.Reader) ([]byte, error)
	Decode(data []byte) (HandlerKey, Message, error)
	Encode(message Message) error
}
