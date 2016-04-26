package ictl

import (
	"sync"
	"sync/atomic"
)

type decoderStats struct {
	recent []bool
	p      int
	mu     *sync.Mutex

	delivered int64
	length    float64
}

func newDecoderStats(length int) (s *decoderStats) {
	s = new(decoderStats)
	s.recent = make([]bool, length)
	for i := range s.recent {
		s.recent[i] = true
	}
	s.mu = new(sync.Mutex)
	s.delivered = int64(length)
	s.length = float64(length)
	return
}

func (s *decoderStats) decoded(successful bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.recent[s.p] != successful {
		s.recent[s.p] = successful
		if successful {
			atomic.AddInt64(&s.delivered, 1)
		} else {
			atomic.AddInt64(&s.delivered, -1)
		}
	}
	s.p++
	if s.p == len(s.recent) {
		s.p = 0
	}
}

func (s *decoderStats) successRatio() float64 {
	return float64(atomic.LoadInt64(&s.delivered)) / s.length
}
