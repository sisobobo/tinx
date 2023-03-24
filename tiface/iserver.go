package tiface

type IServer interface {
	Start()
	Stop()
	SetPack(connect IPack)
	AddRouter(msgId uint32, router IRouter)
}
