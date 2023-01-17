package tnet

import (
	"github.com/sisobobo/tinx/tconf"
	"github.com/sisobobo/tinx/tpkg/bytes"
)

type Round struct {
	readers []bytes.Pool
	writers []bytes.Pool
	buckets []*Bucket
	conf    *tconf.Config
}

func NewRound(c *tconf.Config) (r *Round) {
	r = &Round{
		conf: c,
	}
	// reader
	r.readers = make([]bytes.Pool, c.TCP.Reader)
	for i := 0; i < c.TCP.Reader; i++ {
		r.readers[i].Init(c.TCP.ReadBuf, c.TCP.ReadBufSize)
	}
	// writer
	r.writers = make([]bytes.Pool, c.TCP.Writer)
	for i := 0; i < c.TCP.Writer; i++ {
		r.writers[i].Init(c.TCP.WriteBuf, c.TCP.WriteBufSize)
	}
	//buckets
	r.buckets = make([]*Bucket, c.Bucket.Size)
	for i := 0; i < c.Bucket.Size; i++ {
		r.buckets[i] = NewBucket(c.Bucket)
	}
	return
}

func (r *Round) Reader(rn int) *bytes.Pool {
	return &(r.readers[rn%r.conf.TCP.Reader])
}

func (r *Round) Writer(rn int) *bytes.Pool {
	return &(r.writers[rn%r.conf.TCP.Writer])
}

func (r *Round) Bucket(rn int) *Bucket {
	return r.buckets[rn%r.conf.Bucket.Size]
}
