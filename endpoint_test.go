package ictl

import (
	"bytes"
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
