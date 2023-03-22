package tnet

import (
	"github.com/sisobobo/tinx/tconf"
	"github.com/sisobobo/tinx/tpkg/bytes"
)

type round struct {
	readers []bytes.Pool
	writers []bytes.Pool
}

func newRound(c *tconf.Config) (r *round) {
	var i int
	r = &round{}
	r.readers = make([]bytes.Pool, c.Server.Reader)
	for i = 0; i < c.Server.Reader; i++ {
		r.readers[i].Init(c.Server.ReadBuf, c.Server.ReadBufSize)
	}
	// writer
	r.writers = make([]bytes.Pool, c.Server.Writer)
	for i = 0; i < c.Server.Writer; i++ {
		r.writers[i].Init(c.Server.WriteBuf, c.Server.WriteBufSize)
	}
	return
}

func (r *round) reader(rn int) *bytes.Pool {
	return &(r.readers[rn%len(r.readers)])
}

func (r *round) writer(rn int) *bytes.Pool {
	return &(r.writers[rn%len(r.writers)])
}
