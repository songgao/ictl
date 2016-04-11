package ictl

import (
	"encoding/binary"
	"io"
)

const (
	frameKF uint8 = 1 << iota
	frameDF
)

type CompressionAlgorithm uint8

// Compression algorithms
const (
	CANone CompressionAlgorithm = iota
	CAFlate
	CAGzip
	CALzw
	CAZlib

	// CAAuto is used to indicate auto selecting compression algorithms in
	// encoders. This value is reserved and is never present in ICTL header.
	CAAuto CompressionAlgorithm = 0x0F
)

type header struct {
	frameType          uint8
	compressionOptions uint8
	frameID            uint16
}

func (h *header) setFrameType(frame uint8) {
	// higher (first) 4 bits are reserved and should always be 0b0000
	h.frameType = 0x0F & frame
}

func (h header) getFrameType() uint8 {
	return h.frameType & 0x0F
}

func (h *header) setCompressionOptions(options uint8) {
	// lower 4 bits reserved for compression algorithm
	h.compressionOptions |= 0xF0 & options
}

func (h header) getCompressionOptions() uint8 {
	return h.compressionOptions & 0xF0
}

func (h *header) setCompressionAlgorithm(algo CompressionAlgorithm) {
	// higher 4 bits reserved for parameters
	h.compressionOptions |= 0x0F & uint8(algo)
}

func (h header) getCompressionAlgorithm() CompressionAlgorithm {
	return CompressionAlgorithm(h.compressionOptions & 0x0F)
}

func (h *header) setFrameID(id uint16) {
	h.frameID = id
}

func (h header) getFrameID() uint16 {
	return h.frameID
}

func (h header) writeTo(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, &h.frameType)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, &h.compressionOptions)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, &h.frameID)
	if err != nil {
		return
	}
	return
}

func (h *header) readFrom(r io.Reader) (err error) {
	err = binary.Read(r, binary.BigEndian, &h.frameType)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &h.compressionOptions)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &h.frameID)
	if err != nil {
		return
	}
	return
}
