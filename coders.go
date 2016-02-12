package ictl

import (
	"fmt"
	"sync/atomic"
)

type encoder struct {
	pool    *slicePool
	sentKFs *sliceCacheWithConfidence

	// these two variables are supposed to be uint16; since uint16 is not
	// supported in sync/atomic, we use uint32 here and truncate when using in
	// packet headers
	idCounter   uint32
	cycleLength uint32
}

func newEncoder(pool *slicePool) *encoder {
	return &encoder{
		pool:        pool,
		sentKFs:     newSliceCacheWithConfidence(32),
		cycleLength: 4,
	}
}

// default is 4
func (e *encoder) updateCycleLength(newLength uint16) {
	atomic.StoreUint32(&e.cycleLength, uint32(newLength))
}

func (e *encoder) encode(data *ReusableSlice, confidence uint8) (packet *ReusableSlice, err error) {
	cl := atomic.LoadUint32(&e.cycleLength)
	id := atomic.LoadUint32(&e.idCounter)
	atomic.StoreUint32(&e.idCounter, id+1)

	if id%cl == 0 { // KF; just send the data
		if packet, err = encode(e.pool, data.Slice(), uint16(id), frameKF, CAAuto); err != nil {
			return
		}
		e.sentKFs.put(uint16(id), 0, data) // transferring ownership of data
	} else { // DF; find a proper previously sent KF, and build differential data
		refID, _, ref := e.sentKFs.getMostConfident(4)
		payload := e.pool.get() // temporary buffer to store uncompress data
		defer payload.Done()
		xor(ref.Slice(), data.Slice(), payload)
		data.Done()
		if packet, err = encode(e.pool, payload.Slice(), uint16(refID), frameDF, CAAuto); err != nil {
			return
		}
	}

	return
}

type decoder struct {
	pool    *slicePool
	rcvdKFs *sliceCache
}

func newDecoder(pool *slicePool) *decoder {
	return &decoder{
		pool:    pool,
		rcvdKFs: newSliceCache(32),
	}
}

func (e *decoder) decode(packet []byte) (data *ReusableSlice, err error) {
	var header header
	var payload *ReusableSlice
	if header, payload /* uncompressed payload */, err = decode(e.pool, packet); err != nil {
		return
	}

	if header.frameType == frameKF { // in KF, uncompressed payload is the data
		payload.AddOwner()
		e.rcvdKFs.put(header.frameID, payload) // 1st owner
		data = payload                         // 2nd owner
	} else if header.frameType == frameDF { // in DF, uncompressed payload is differential data
		defer payload.Done()
		ref, ok := e.rcvdKFs.get(header.frameID)
		if !ok {
			err = fmt.Errorf("referenced frame (id=%d)is missing", header.frameID)
		}
		data = e.pool.get()
		xor(ref.Slice(), payload.Slice(), data)
	}

	return
}
