package tnet

import "tinx/tiface"

type MsgHandler struct {
}

func (m *MsgHandler) PreHandler(channel tiface.IChannel, message interface{}) {
}

func (m *MsgHandler) Handler(channel tiface.IChannel, message interface{}) {
}

func (m *MsgHandler) PostHandler(channel tiface.IChannel, message interface{}) {

}
