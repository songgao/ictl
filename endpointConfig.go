package ictl

type EndpointConfig interface {
	MaxPacketSize() int
	CompressionAlgorithm() compressionAlgorithm

	SetMaxPacketSize(int) EndpointConfig
	SetCompressionAlgorithm(compressionAlgorithm) EndpointConfig
}

func DefaultEndpointConfig() EndpointConfig {
	return &endpointConfig{
		maxPacketSize: 1379,
		cmpAlgr:       CAAuto,
	}
}

type endpointConfig struct {
	maxPacketSize int
	cmpAlgr       compressionAlgorithm
}

func (e *endpointConfig) MaxPacketSize() int                         { return e.maxPacketSize }
func (e *endpointConfig) CompressionAlgorithm() compressionAlgorithm { return e.cmpAlgr }

func (e *endpointConfig) SetMaxPacketSize(v int) EndpointConfig {
	e.maxPacketSize = v
	return e
}

func (e *endpointConfig) SetCompressionAlgorithm(v compressionAlgorithm) EndpointConfig {
	e.cmpAlgr = v
	return e
}
