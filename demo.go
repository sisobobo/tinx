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

func (t *TestMsg) RouterId() any {
	return nil
}

func NewTestMsg(s string) *TestMsg {
	m := new(TestMsg)
	m.msg = s
	return m
}

type TestCodec struct {
}

func (t *TestCodec) Decode(reader *bufio.Reader) (tnet.Message, error) {
	line, err := reader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	s := string(line[:len(line)-2])
	msg := NewTestMsg(s)
	return msg, err
}

func (t *TestCodec) Encode(message tnet.Message) ([]byte, error) {
	return nil, nil
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

func aa() {
	fmt.Println("11")
	return
}

func bb() {
	fmt.Println("22")
}

func main() {
	//server := tnet.NewServer("",
	//	tnet.SetCodec(&TestCodec{}),
	//	tnet.SetHandler(TestHandler{}),
	//)
	//server.Serve()
}
