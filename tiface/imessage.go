package tiface

type IMessage interface {
	GetMsgId() uint32
	GetData() []byte
}
