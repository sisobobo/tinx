package tconf

type Config struct {
	Debug  bool
	Env    *Env
	Bucket *Bucket
	Server *Server
}

type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
	Weight    int64
	Offline   bool
	Addrs     []string
}

type Bucket struct {
	Size    int
	Channel int
}

type Server struct {
	Bind         []string
	SndBuf       int
	RcvBuf       int
	KeepAlive    bool
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}
