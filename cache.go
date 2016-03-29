package ictl

import "container/ring"

func initRing(r *ring.Ring, value interface{}) *ring.Ring {
	r.Value = value
	for i := r.Next(); i != nil && i != r; i = i.Next() {
		i.Value = value
	}
	return r
}

type sliceCache struct {
	lastID *ring.Ring
	slices map[uint16]*ReusableSlice
}

func newSliceCache(size int) *sliceCache {
	return &sliceCache{
		lastID: initRing(ring.New(size), uint16(0)),
		slices: make(map[uint16]*ReusableSlice),
	}
}

func (c *sliceCache) put(id uint16, slice *ReusableSlice) {
	c.lastID = c.lastID.Next()
	oldId := c.lastID.Value.(uint16)
	if oldSlice, ok := c.slices[oldId]; ok {
		oldSlice.Done()
		delete(c.slices, oldId)
	}
	c.lastID.Value = id
	c.slices[id] = slice
}

func (c *sliceCache) get(id uint16) (slice *ReusableSlice, ok bool) {
	if slice, ok = c.slices[id]; ok {
		slice.AddOwner()
	}
	return
}

type sliceWithConfidence struct {
	slice      *ReusableSlice
	confidence uint8
}
type sliceCacheWithConfidence struct {
	lastID *ring.Ring
	slices map[uint16]sliceWithConfidence
}

func newSliceCacheWithConfidence(size int) (c *sliceCacheWithConfidence) {
	return &sliceCacheWithConfidence{
		lastID: initRing(ring.New(size), uint16(0)),
		slices: make(map[uint16]sliceWithConfidence),
	}
}

func (c *sliceCacheWithConfidence) put(id uint16, confidence uint8, slice *ReusableSlice) {
	c.lastID = c.lastID.Next()
	oldId := c.lastID.Value.(uint16)
	if oldSlice, ok := c.slices[oldId]; ok {
		oldSlice.slice.Done()
		delete(c.slices, oldId)
	}
	c.lastID.Value = id
	c.slices[id] = sliceWithConfidence{slice: slice, confidence: confidence}
}

// Get the slice with largest confidence value, within last num slices inserted
// by Put()
func (c *sliceCacheWithConfidence) getMostConfident(num int) (id uint16, confidence uint8, slice *ReusableSlice) {
	if num <= 0 {
		panic(nil)
	}

	r := c.lastID
	for i := num; i > 0; i-- {
		currentID := r.Value.(uint16)
		if currentSlice, currentOK := c.slices[currentID]; i == num || (currentOK && currentSlice.confidence > confidence) {
			id = currentID
			confidence = currentSlice.confidence
			slice = currentSlice.slice
		}
		r = r.Prev()
	}

	slice.AddOwner()

	return
}
