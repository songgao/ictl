package ictl

import (
	"bytes"
	"errors"
	"io"
)

func compress(c compressor, output []byte, data []byte) (length int, err error) {
	buf := bytes.NewBuffer(output[0:0:cap(output)])
	var w io.WriteCloser
	if w, err = c.compressor(buf); err != nil {
		return
	}
	if _, err = w.Write(data); err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}
	length = buf.Len()
	return
}

func compressFindBest(output []byte, data []byte) (cmp compressor, length int, err error) {
	var exhaustive []compressor
	for _, v := range compressors {
		exhaustive = append(exhaustive, v())
	}

	var l int
	best, bestk := int((^uint(0))>>1), 255
	for k, c := range exhaustive {
		if l, err = compress(c, output, data); err != nil {
			return
		}
		if l < best {
			bestk = k
			best = l
		}
	}

	cmp = exhaustive[bestk]
	if length, err = compress(cmp, output, data); err != nil {
		return
	}

	return
}

func encode(pool *slicePool, payload []byte, id uint16, frameType uint8, cmpAlgr compressionAlgorithm) (packet *ReusableSlice, err error) {
	packet = pool.get()
	cleanup := func() {
		packet.Done()
		packet = nil
	}

	var cmp compressor
	var l int
	if cmpAlgr == CAAuto {
		if cmp, l, err = compressFindBest(packet.Slice()[4:], payload); err != nil {
			cleanup()
			return
		}
	} else {
		var creator compressorCreator
		var ok bool
		if creator, ok = compressors[cmpAlgr]; !ok {
			cleanup()
			err = errors.New("unknown compression algorithm")
			return
		}
		cmp = creator()
		if l, err = compress(cmp, packet.Slice()[4:], payload); err != nil {
			cleanup()
			return
		}
	}
	packet.Resize(l + 4)

	var header header
	header.setFrameID(id)
	header.setFrameType(frameType)
	header.setCompressionOptions(cmp.getOptionsForHeader())
	header.setCompressionAlgorithm(cmp.getCompressionAlgorithm())
	if err = header.writeTo(bytes.NewBuffer(packet.Slice()[0:0:4])); err != nil {
		cleanup()
		return
	}

	return
}

func decode(pool *slicePool, packet []byte) (header header, payload *ReusableSlice, err error) {
	payload = pool.get()
	cleanup := func() {
		payload.Done()
		payload = nil
	}

	reader := bytes.NewReader(packet)
	header.readFrom(reader)
	c := compressors[header.getCompressionAlgorithm()]()
	c.setOptionsFromHeader(header.getCompressionOptions())
	var r io.ReadCloser
	if r, err = c.decompressor(reader); err != nil {
		cleanup()
		return
	}
	var l int
	if l, err = r.Read(payload.Slice()); err != nil {
		cleanup()
		return
	}
	if err = r.Close(); err != nil {
		cleanup()
		return
	}
	payload.Resize(l)

	return
}

// calling xor, unlike other functions, doesn't transfer ownership
func xor(a, b []byte, output *ReusableSlice) {
	c := output.Slice()
	var i int
	lastNonZero := -1
	for i = 0; i < len(a) && i < len(b); i++ {
		c[i] = a[i] ^ b[i]
		if c[i] != 0 {
			lastNonZero = i
		}
	}
	for ; i < len(a); i++ {
		c[i] = a[i] ^ 0
		if c[i] != 0 {
			lastNonZero = i
		}
	}
	for ; i < len(b); i++ {
		c[i] = b[i] ^ 0
		if c[i] != 0 {
			lastNonZero = i
		}
	}
	output.Resize(lastNonZero + 1)
}
