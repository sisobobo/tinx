package tiface

type IServer interface {
	Start()
	Stop()
	SetPack(pack IPack)
	AddMsgHandlers(handlers ...IHandler)
}
