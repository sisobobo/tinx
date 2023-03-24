package tiface

type IMessage interface {
	GetMsgId() uint32
	GetData() []byte
	SetMsgId(id uint32)
	SetData(data []byte)
}
