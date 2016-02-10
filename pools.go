package ictl

import "sync"

type ReusableSlice struct {
	slice []byte
	pool  *sync.Pool
}

func (s *ReusableSlice) Done() {
	s.pool.Put(s)
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

func (p *slicePool) Get() *ReusableSlice {
	s := p.pool.Get().(*ReusableSlice)
	s.Resize(s.Cap())
	return s
}
