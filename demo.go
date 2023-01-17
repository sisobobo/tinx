package main

import (
	"fmt"
	"github.com/sisobobo/tinx/tnet"
	"github.com/sisobobo/tinx/tpkg/bufio"
)

type TestMsg struct {
	tnet.Message
	msg string
}

func NewTestMsg(s string) *TestMsg {
	m := new(TestMsg)
	m.msg = s
	return m
}

func (m *TestMsg) String() string {
	return m.msg
}

type TestCodec struct {
}

func (t *TestCodec) Decode(reader *bufio.Reader) (tnet.Message, error) {
	line, err := reader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	s := string(line[:len(line)-2])
	return s, err
}

func (t *TestCodec) Encode() {

}

type TestHandler struct {
}

func (t TestHandler) Connect(channel *tnet.Channel) {
	fmt.Println(channel.RemoteAddr(), "已连接")
}

func (t TestHandler) Receive(channel *tnet.Channel, message tnet.Message) {
	fmt.Println(channel.RemoteAddr(), ":", message)
}

func (t TestHandler) DisConnect(channel *tnet.Channel) {
	fmt.Println(channel.RemoteAddr(), "断开连接")
}

func main() {
	server := tnet.NewServer("",
		tnet.SetCodec(&TestCodec{}),
		tnet.SetHandler(TestHandler{}),
	)
	server.Serve()
}
