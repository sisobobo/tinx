package tiface

type Server interface {
	Start()
	Stop()
	Serve()
}
