package ictl

import "sync"

type Endpoint interface {
	Encode(context string, data []byte, confidence uint8) (packet *ReusableSlice, err error)
	EncodeReusable(context string, data *ReusableSlice, confidence uint8) (packet *ReusableSlice, err error)
	Decode(context string, packet []byte) (data *ReusableSlice, err error)
}

type endpoint struct {
	config endpointConfig
	pool   *slicePool

	encoders map[string]*encoder
	decoders map[string]*decoder
	mapMu    *sync.Mutex

	dStats *decoderStats
}

func NewEndpoint(config EndpointConfig) Endpoint {
	e := &endpoint{
		config: *(config.(*endpointConfig)), // copy
	}
	e.pool = newSlicePool(e.config.maxPacketSize)
	e.encoders = make(map[string]*encoder)
	e.decoders = make(map[string]*decoder)
	e.mapMu = new(sync.Mutex)

	e.dStats = newDecoderStats(100)

	return e
}

// data is copied to a ReusableSlice internally, i.e., caller can use the data
// slice for other purposes safely
func (e *endpoint) Encode(context string, data []byte, confidence uint8) (packet *ReusableSlice, err error) {
	d := e.pool.get()
	copy(d.Slice(), data)
	d.Resize(len(data))
	packet, err = e.EncodeReusable(context, d, confidence)
	return
}

func (e *endpoint) EncodeReusable(context string, data *ReusableSlice, confidence uint8) (packet *ReusableSlice, err error) {
	e.mapMu.Lock()
	enc, ok := e.encoders[context]
	if !ok {
		enc = newEncoder(e.pool, e.config, e.dStats)
		e.encoders[context] = enc
	}
	e.mapMu.Unlock()
	packet, err = enc.encode(data, confidence)
	return
}

func (e *endpoint) Decode(context string, packet []byte) (data *ReusableSlice, err error) {
	e.mapMu.Lock()
	dec, ok := e.decoders[context]
	if !ok {
		dec = newDecoder(e.pool, e.dStats)
		e.decoders[context] = dec
	}
	e.mapMu.Unlock()
	data, err = dec.decode(packet)
	return
}
