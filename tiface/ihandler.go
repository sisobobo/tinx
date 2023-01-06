package tiface

type IHandler interface {
	Id() interface{}                                   //此handler唯一标识
	PreHandler(channel IChannel, message interface{})  //在处理conn业务之前的钩子方法
	Handler(channel IChannel, message interface{})     //处理conn业务的方法
	PostHandler(channel IChannel, message interface{}) //处理conn业务之后的钩子方法
}
