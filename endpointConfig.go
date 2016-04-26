package ictl

type EndpointConfig interface {
	MaxPacketSize() int
	CompressionAlgorithm() CompressionAlgorithm
	EncoderCycleLength() uint16
	ConfidenceLookback() int

	SetConfidenceLookback(int) EndpointConfig
	SetMaxPacketSize(int) EndpointConfig
	SetCompressionAlgorithm(CompressionAlgorithm) EndpointConfig

	// set to 0 to use adaptive
	SetEncoderCycleLength(uint16) EndpointConfig
}

func DefaultEndpointConfig() EndpointConfig {
	return &endpointConfig{
		maxPacketSize:      1379,
		cmpAlgr:            CAAuto,
		cycleLength:        0,
		confidenceLookback: 1,
	}
}

type endpointConfig struct {
	maxPacketSize      int
	cmpAlgr            CompressionAlgorithm
	cycleLength        uint16
	confidenceLookback int
}

func (e *endpointConfig) MaxPacketSize() int                         { return e.maxPacketSize }
func (e *endpointConfig) CompressionAlgorithm() CompressionAlgorithm { return e.cmpAlgr }
func (e *endpointConfig) EncoderCycleLength() uint16                 { return e.cycleLength }
func (e *endpointConfig) ConfidenceLookback() int                    { return e.confidenceLookback }

func (e *endpointConfig) SetMaxPacketSize(v int) EndpointConfig {
	e.maxPacketSize = v
	return e
}

func (e *endpointConfig) SetCompressionAlgorithm(v CompressionAlgorithm) EndpointConfig {
	e.cmpAlgr = v
	return e
}

func (e *endpointConfig) SetEncoderCycleLength(cycleLength uint16) EndpointConfig {
	e.cycleLength = cycleLength
	return e
}

func (e *endpointConfig) SetConfidenceLookback(confidenceLookback int) EndpointConfig {
	e.confidenceLookback = confidenceLookback
	if confidenceLookback < 1 {
		panic("invalid confidenceLoopback")
	}
	return e
}
