package tnet

type Handler interface {
	Connect(channel *Channel)
	Receive(channel *Channel, msg Message)
	DisConnect(channel *Channel)
}

type RouterHandler struct {
	Handler
}

func (rh *RouterHandler) router(rm *RouterManager, channel *Channel, msg Message) {
	rh.Receive(channel, msg)
	if rm != nil {
		rm.doHandler(channel, msg)
	}
}

func NewRouterHandler(handler Handler) *RouterHandler {
	return &RouterHandler{Handler: handler}
}
