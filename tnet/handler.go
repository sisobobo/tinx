package tnet

type Handler interface {
	Connect(channel *Channel)
	Receive(channel *Channel, msg Message)
	DisConnect(channel *Channel)
}
