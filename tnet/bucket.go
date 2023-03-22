package tnet

import (
	"github.com/sisobobo/tinx/tconf"
	"sync"
)

type bucket struct {
	c     *tconf.Bucket
	cLock sync.RWMutex
	chs   map[string]*channel
}

func newBucket(c *tconf.Bucket) *bucket {
	return &bucket{
		c:   c,
		chs: make(map[string]*channel, c.Channel),
	}
}

func (b *bucket) putChannel(ch *channel) {
	b.cLock.Lock()
	defer b.cLock.Unlock()
	if dch := b.chs[ch.id]; dch != nil {
		dch.Close()
	}
	b.chs[ch.id] = ch
}

func (b *bucket) delChannel(dch *channel) {
	b.cLock.Lock()
	defer b.cLock.Unlock()
	if ch, ok := b.chs[dch.id]; ok {
		if ch == dch {
			delete(b.chs, ch.id)
		}
		b.chs[ch.id] = ch
	}
}
