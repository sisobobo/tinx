package tconf

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Debug  bool
	TCP    *TCP
	Bucket *Bucket
}

type TCP struct {
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

type Bucket struct {
	Size          int
	Channel       int
	Room          int
	RoutineAmount uint64
	RoutineSize   int
}

func NewConfig(confPath string) (conf *Config, err error) {
	conf = Default()
	if len(confPath) > 0 {
		_, err = toml.DecodeFile(confPath, &conf)
	}
	return conf, err
}

func Default() *Config {
	return &Config{
		Debug: true,
		TCP: &TCP{
			Bind:         []string{":7777"},
			SndBuf:       4096,
			RcvBuf:       4096,
			KeepAlive:    false,
			Reader:       32,
			ReadBuf:      1024,
			ReadBufSize:  8192,
			Writer:       32,
			WriteBuf:     1024,
			WriteBufSize: 8192,
		},
		Bucket: &Bucket{
			Size:          32,
			Channel:       1024,
			Room:          1024,
			RoutineAmount: 32,
			RoutineSize:   1024,
		},
	}
}
