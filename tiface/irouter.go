package tiface

type IRouter interface {
	PreHandle(channel IChannel, msg IMessage)
	Handle(channel IChannel, msg IMessage)
	PostHandle(channel IChannel, msg IMessage)
}
