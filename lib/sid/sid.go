package sid

import (
	"sync"
	"time"
)

type Sid struct {
	id    uint64
	Index uint64
	lock  sync.Mutex
}

func (v *Sid) SetIndex(i uint) {

	v.Index = uint64(i)
}

func (v *Sid) GetId() uint64 {
	tnow := time.Now().Unix()

	var s uint64
	s = uint64(tnow)
	s *= 100
	s += v.Index
	s *= 1e7

	v.lock.Lock()
	v.id += 1
	if v.id > 1e6 {
		v.id = 1
	}
	s += v.id
	v.lock.Unlock()
	return s
}
