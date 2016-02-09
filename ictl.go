package ictl

type EndpointConfig struct {
	MaxPayloadSize int
}

func DefaultEndpointConfig() *EndpointConfig {
	return &EndpointConfig{
		MaxPayloadSize: 1395,
	}
}

type Endpoint struct {
	config *EndpointConfig
}

func NewEndpoint(config *EndpointConfig) *Endpoint {
	e := &Endpoint{
		config: config,
	}
	return e
}
