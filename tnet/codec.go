package tnet

import "github.com/sisobobo/tinx/tpkg/bufio"

type Message interface {
	RouterId() interface{}
}

type Codec interface {
	Decode(*bufio.Reader) (Message, error)
	Encode(Message) ([]byte, error)
}
