package ictl

import "fmt"

type encoder struct {
	pool    *slicePool
	sentKFs *sliceCacheWithConfidence

	idCounter          uint16
	cycleLength        uint16
	confidenceLookback int
	cmpAlgr            CompressionAlgorithm

	adaptive *adaptiveCycleLength
	dStats   *decoderStats
}

func newEncoder(pool *slicePool, config endpointConfig, dStats *decoderStats) *encoder {
	return &encoder{
		pool:               pool,
		sentKFs:            newSliceCacheWithConfidence(32),
		cycleLength:        config.cycleLength,
		confidenceLookback: config.confidenceLookback,
		cmpAlgr:            config.cmpAlgr,
		adaptive:           new(adaptiveCycleLength),
		dStats:             dStats,
	}
}

func (e *encoder) encKF(id uint16, data *ReusableSlice, confidence uint8) (packet *ReusableSlice, err error) {
	if packet, err = encode(e.pool, data.Slice(), uint16(id), frameKF, e.cmpAlgr); err != nil {
		return
	}
	e.sentKFs.put(uint16(id), confidence, data) // transferring ownership of data
	return
}

func (e *encoder) encDF(data *ReusableSlice, confidence uint8) (packet *ReusableSlice, err error) {
	refID, _, ref := e.sentKFs.getMostConfident(e.confidenceLookback)
	payload := e.pool.get() // temporary buffer to store uncompress data
	defer payload.Done()
	xor(ref.Slice(), data.Slice(), payload)
	data.Done()
	if packet, err = encode(e.pool, payload.Slice(), uint16(refID), frameDF, e.cmpAlgr); err != nil {
		return
	}
	return
}

func (e *encoder) encode(data *ReusableSlice, confidence uint8) (packet *ReusableSlice, err error) {
	if e.cycleLength != 0 { // fixed cycle length
		if e.idCounter%e.cycleLength == 0 { // KF; just send the data
			packet, err = e.encKF(e.idCounter, data, confidence)
		} else { // DF; find a proper previously sent KF, and build differential data
			packet, err = e.encDF(data, confidence)
		}
	} else { // adaptive cycle length
		if e.adaptive.first() {
			packet, err = e.encKF(e.idCounter, data, confidence)
			e.adaptive.sentKF(len(packet.Slice()))
		} else {
			data.AddOwner()
			packet, err = e.encDF(data, confidence)
			if e.adaptive.shouldSendThisDF(len(packet.Slice())) {
				e.adaptive.sentDF(len(packet.Slice()))
				data.Done()
			} else {
				packet.Done()
				packet, err = e.encKF(e.idCounter, data, confidence)
				e.adaptive.sentKF(len(packet.Slice()))
			}
		}
	}

	e.idCounter++

	return
}

type decoder struct {
	pool    *slicePool
	rcvdKFs *sliceCache

	dStats *decoderStats
}

func newDecoder(pool *slicePool, dStats *decoderStats) *decoder {
	return &decoder{
		pool:    pool,
		rcvdKFs: newSliceCache(32),
		dStats:  dStats,
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
		e.dStats.decoded(ok)
		if !ok {
			err = fmt.Errorf("referenced frame (id=%d)is missing", header.frameID)
			return
		}
		data = e.pool.get()
		xor(ref.Slice(), payload.Slice(), data)
	}

	return
}
