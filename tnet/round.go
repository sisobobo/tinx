package tnet

import (
	"github.com/sisobobo/tinx/tpkg/bytes"
	"time"
)

type RoundOptions struct {
	Timer        int
	TimerSize    int
	Readers      int
	ReadBuf      int
	ReadBufSize  int
	Writers      int
	WriteBuf     int
	WriteBufSize int
}

func defaultRoundOptions() *RoundOptions {
	return &RoundOptions{
		Readers:      32,
		ReadBuf:      1024,
		ReadBufSize:  8192,
		Writers:      32,
		WriteBuf:     1024,
		WriteBufSize: 8192,
		Timer:        32,
		TimerSize:    2048,
	}
}

type Round struct {
	readers []bytes.Pool
	writers []bytes.Pool
	timers  []time.Timer
	options *RoundOptions
}

func newRound() *Round {
	r := &Round{
		options: defaultRoundOptions(),
	}
	r.readers = make([]bytes.Pool, r.options.Readers)
	for i := 0; i < r.options.Readers; i++ {
		r.readers[i].Init(r.options.ReadBuf, r.options.ReadBufSize)
	}
	r.writers = make([]bytes.Pool, r.options.Writers)
	for i := 0; i < r.options.Writers; i++ {
		r.writers[i].Init(r.options.WriteBuf, r.options.WriteBufSize)
	}
	return r
}

func (r *Round) Reader(n int) *bytes.Pool {
	return &r.readers[n%r.options.Readers]
}

func (r *Round) Writer(n int) *bytes.Pool {
	return &r.readers[n%r.options.Readers]
}

func (r *Round) Timer(n int) *time.Timer {
	return &(r.timers[n%r.options.Timer])
}
