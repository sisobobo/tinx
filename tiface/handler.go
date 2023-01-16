package tiface

type Handler interface {
	Active(channel Channel)
	Read(channel Channel, data interface{})
	InActive(channel Channel)
}
