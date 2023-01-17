package tnet

import (
	"github.com/sisobobo/tinx/tconf"
	"sync"
)

type Bucket struct {
	conf  *tconf.Bucket
	cLock sync.RWMutex
	chs   map[string]*Channel
}

func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	defer b.cLock.RUnlock()
	ch = b.chs[key]
	return
}

func (b *Bucket) Put(ch *Channel) (err error) {
	b.cLock.Lock()
	defer b.cLock.Unlock()
	if dch := b.chs[ch.id]; dch != nil {
		dch.Close()
	}
	b.chs[ch.id] = ch
	return
}

func (b *Bucket) Remove(ch *Channel) {
	b.cLock.Lock()
	defer b.cLock.Unlock()
	if v, ok := b.chs[ch.id]; ok {
		if v == ch {
			delete(b.chs, ch.id)
		}
	}
}

func NewBucket(c *tconf.Bucket) *Bucket {
	b := &Bucket{
		conf: c,
		chs:  make(map[string]*Channel, c.Channel),
	}
	return b
}
