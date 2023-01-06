package tiface

type IMessage interface {
	HandlerId() interface{}
	Msg() interface{}
}
