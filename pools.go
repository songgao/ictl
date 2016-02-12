package ictl

import "sync"
import "sync/atomic"

type ReusableSlice struct {
	slice   []byte
	pool    *sync.Pool
	counter int32
}

func (s *ReusableSlice) AddOwner() {
	atomic.AddInt32(&s.counter, 1)
}

func (s *ReusableSlice) Done() {
	atomic.AddInt32(&s.counter, -1)
	c := atomic.LoadInt32(&s.counter)
	if c == 0 {
		s.pool.Put(s)
	} else if c < 0 {
		panic("incorrect use of ReusableSlice")
	}
}

func (s *ReusableSlice) Slice() []byte {
	return s.slice
}

// Shrink or grow the slice; length cannot exceed Cap()
func (s *ReusableSlice) Resize(length int) {
	s.slice = s.slice[0:length]
}

func (s *ReusableSlice) Cap() int {
	return cap(s.slice)
}

type slicePool struct {
	pool *sync.Pool
}

func newSlicePool(maxLength int) *slicePool {
	pool := new(sync.Pool)
	pool.New = func() interface{} {
		s := new(ReusableSlice)
		s.pool = pool
		s.slice = make([]byte, maxLength)
		return s
	}
	return &slicePool{
		pool: pool,
	}
}

func (p *slicePool) get() *ReusableSlice {
	s := p.pool.Get().(*ReusableSlice)
	s.Resize(s.Cap())
	s.counter = 1
	return s
}
