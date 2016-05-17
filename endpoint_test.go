package ictl

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"
)

func TestEndpoint(t *testing.T) {
	config := DefaultEndpointConfig()
	endpoint1 := NewEndpoint(config)
	endpoint2 := NewEndpoint(config)

	var packet, rcvd *ReusableSlice
	var err error
	for i := 0; i < 10; i++ {
		var toSend []byte
		if toSend, err = time.Now().MarshalBinary(); err != nil {
			t.Fatalf("calling time.Now().MarshalBinary() error: %v\n", err)
		}
		if packet, err = endpoint1.Encode("test", toSend, 0); err != nil {
			t.Fatalf("calling endpoint1.Encode() error: %v\n", err)
		}
		if rcvd, err = endpoint2.Decode("test", packet.Slice()); err != nil {
			t.Fatalf("calling endpoint2.Decode() error: %v\n", err)
		}
		packet.Done()
		if !bytes.Equal(toSend, rcvd.Slice()) {
			t.Fatalf("decoded data is not equal to sent data: %v != %v\n", toSend, rcvd.Slice())
		} else {
			t.Logf("encoding/decoding passed for: %x\n", toSend)
		}
		rcvd.Done()
	}
}

func BenchmarkDefaultRandom(b *testing.B) {
	config := DefaultEndpointConfig()
	endpoint := NewEndpoint(config)
	data := make([][]byte, 256)
	for i := range data {
		data[i] = make([]byte, 256)
		if _, err := rand.Read(data[i]); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if p, err := endpoint.Encode("test", data[i%256], 0); err != nil {
			b.Fatal(err)
		} else {
			p.Done()
		}
	}
}

func BenchmarkDefaultFlateRandom(b *testing.B) {
	config := DefaultEndpointConfig().SetCompressionAlgorithm(CAFlate)
	endpoint := NewEndpoint(config)
	data := make([][]byte, 256)
	for i := range data {
		data[i] = make([]byte, 256)
		if _, err := rand.Read(data[i]); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if p, err := endpoint.Encode("test", data[i%256], 0); err != nil {
			b.Fatal(err)
		} else {
			p.Done()
		}
	}
}

func BenchmarkDefaultIdentical(b *testing.B) {
	config := DefaultEndpointConfig()
	endpoint := NewEndpoint(config)
	data := make([]byte, 256)
	if _, err := rand.Read(data); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if p, err := endpoint.Encode("test", data, 0); err != nil {
			b.Fatal(err)
		} else {
			p.Done()
		}
	}
}

func BenchmarkDefaultFlateIdentical(b *testing.B) {
	config := DefaultEndpointConfig().SetCompressionAlgorithm(CAFlate)
	endpoint := NewEndpoint(config)
	data := make([]byte, 256)
	if _, err := rand.Read(data); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if p, err := endpoint.Encode("test", data, 0); err != nil {
			b.Fatal(err)
		} else {
			p.Done()
		}
	}
}
