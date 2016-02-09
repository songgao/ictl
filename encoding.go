package ictl

import (
	"bytes"
	"io"
)

func compressFindBest(data []byte) (algr uint8, options uint8, buf *ReusableBuffer, err error) {
	var exhaustive []struct {
		c compressor
		b *ReusableBuffer
	}
	for _, v := range compressors {
		exhaustive = append(exhaustive, struct {
			c compressor
			b *ReusableBuffer
		}{
			c: v,
			b: poolBuffer.Get(),
		})
	}

	best, besti := int((^uint(0))>>1), -1
	for i, c := range exhaustive {
		var w io.WriteCloser
		if w, err = c.c.compressor(c.b); err != nil {
			return
		}
		if _, err = w.Write(data); err != nil {
			return
		}
		if err = w.Close(); err != nil {
			return
		}
		if c.b.Len() < best {
			besti = i
		}
	}

	for i, c := range exhaustive {
		if i != besti {
			c.b.Done()
		}
	}

	buf = exhaustive[besti].b
	algr = exhaustive[besti].c.getCompressionAlgorithm()
	options = exhaustive[besti].c.getOptionsForHeader()

	return
}

func encodeKF(data []byte, id uint16) (header header, payload *ReusableBuffer, err error) {
	header.setFrameID(id)
	header.setFrameType(frameKF)
	var algr, options uint8
	if algr, options, payload, err = compressFindBest(data); err != nil {
		return
	}
	header.setCompressionOptions(options)
	header.setCompressionAlgorithm(algr)
	return
}

func xor(a, b []byte, output *ReusableSlice) {
	c := output.Slice()
	var i, lastNonZero int
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

func encodeDF(data []byte, base []byte, pool *slicePool, baseId uint16) (header header, payload *ReusableBuffer, err error) {
	header.setFrameID(baseId)
	header.setFrameType(frameDF)
	slice := pool.Get()
	xor(data, base, slice)
	var algr, options uint8
	algr, options, payload, err = compressFindBest(slice.Slice())
	slice.Done()
	if err != nil {
		return
	}
	header.setCompressionOptions(options)
	header.setCompressionAlgorithm(algr)
	return
}

func decode(packet []byte) (header header, data *ReusableBuffer, err error) {
	reader := bytes.NewReader(packet)
	header.readFrom(reader)
	var r io.ReadCloser
	if r, err = compressors[header.getCompressionAlgorithm()].decompressor(bytes.NewReader(packet)); err != nil {
		return
	}
	data = poolBuffer.Get()
	if _, err = data.ReadFrom(r); err != nil {
		return
	}
	if err = r.Close(); err != nil {
		return
	}
	return
}
