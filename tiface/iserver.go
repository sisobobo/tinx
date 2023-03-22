package tiface

type IServer interface {
	Start()
	Stop()
	AddRouter(msgId uint32, router IRouter)
	SetConnect(connect IConnect)
}
