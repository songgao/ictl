package ictl

import "log"

type EndpointConfig interface {
	MaxPacketSize() int
	CompressionAlgorithm() compressionAlgorithm
	EncoderCycleLength() uint32

	SetMaxPacketSize(int) EndpointConfig
	SetCompressionAlgorithm(compressionAlgorithm) EndpointConfig
	SetEncoderCycleLength(uint32) EndpointConfig
}

func DefaultEndpointConfig() EndpointConfig {
	return &endpointConfig{
		maxPacketSize: 1379,
		cmpAlgr:       CAAuto,
		cycleLength:   4,
	}
}

type endpointConfig struct {
	maxPacketSize int
	cmpAlgr       compressionAlgorithm
	cycleLength   uint32
}

func (e *endpointConfig) MaxPacketSize() int                         { return e.maxPacketSize }
func (e *endpointConfig) CompressionAlgorithm() compressionAlgorithm { return e.cmpAlgr }
func (e *endpointConfig) EncoderCycleLength() uint32                 { return e.cycleLength }

func (e *endpointConfig) SetMaxPacketSize(v int) EndpointConfig {
	e.maxPacketSize = v
	return e
}

func (e *endpointConfig) SetCompressionAlgorithm(v compressionAlgorithm) EndpointConfig {
	e.cmpAlgr = v
	return e
}

func (e *endpointConfig) SetEncoderCycleLength(cycleLength uint32) EndpointConfig {
	if cycleLength > 0 {
		e.cycleLength = cycleLength
	} else {
		log.Printf("invalid cycle length %d\n", cycleLength)
	}
	return e
}
