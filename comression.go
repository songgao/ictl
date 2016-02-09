package ictl

import (
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"io"
	"io/ioutil"
)

var compressors map[uint8]compressor = map[uint8]compressor{
	cmpAlgrNone:  compressorNone{},
	cmpAlgrFlate: compressorFlate{},
	cmpAlgrGzip:  compressorGzip{},
	cmpAlgrLzw:   compressorLzw{},
	cmpAlgrZlib:  compressorZlib{},
}

type compressor interface {
	// compressor returns a new WriteCloser that can be used to write
	// uncompressed data, which will be compressed and written into compressed.
	compressor(compressed io.Writer) (uncompressed io.WriteCloser, err error)

	// see `compress/*` packages docs;
	// decompressor returns a new ReadCloser that can be used to read the
	// uncompressed version of compressed. If compressed does not also implement
	// io.ByteReader, the decompressor may read more data than necessary from
	// compressed.  It is the caller's responsibility to call Close on the
	// ReadCloser when finished reading.
	decompressor(compressed io.Reader) (uncompressed io.ReadCloser, err error)
	getOptionsForHeader() uint8     // only higher 4 bits
	getCompressionAlgorithm() uint8 // only lower 4 bits
}

type emptyCompressorOptions struct{}

func (e emptyCompressorOptions) getOptionsForHeader() uint8 {
	return 0
}

type compressorNone struct {
	emptyCompressorOptions
}

func (c compressorNone) compressor(compressed io.Writer) (w io.WriteCloser, err error) {
	w = newWriterNopCloser(compressed)
	return
}

func (c compressorNone) decompressor(compressed io.Reader) (r io.ReadCloser, err error) {
	r = ioutil.NopCloser(compressed)
	return
}

func (c compressorNone) getCompressionAlgorithm() uint8 {
	return cmpAlgrNone
}

type compressorFlate struct {
	emptyCompressorOptions
}

func (c compressorFlate) compressor(compressed io.Writer) (w io.WriteCloser, err error) {
	if w, err = flate.NewWriter(w, 9); err != nil {
		return
	}
	return
}

func (c compressorFlate) decompressor(compressed io.Reader) (r io.ReadCloser, err error) {
	r = flate.NewReader(compressed)
	return
}

func (c compressorFlate) getCompressionAlgorithm() uint8 {
	return cmpAlgrFlate
}

type compressorGzip struct {
	emptyCompressorOptions
}

func (c compressorGzip) compressor(compressed io.Writer) (w io.WriteCloser, err error) {
	w = gzip.NewWriter(w)
	return
}

func (c compressorGzip) decompressor(compressed io.Reader) (r io.ReadCloser, err error) {
	r, err = gzip.NewReader(r)
	if err != nil {
		return
	}
	return
}

func (c compressorGzip) getCompressionAlgorithm() uint8 {
	return cmpAlgrGzip
}

type compressorLzw struct {
	emptyCompressorOptions
}

func (c compressorLzw) compressor(compressed io.Writer) (w io.WriteCloser, err error) {
	w = lzw.NewWriter(w, lzw.MSB, 8)
	return
}

func (c compressorLzw) decompressor(compressed io.Reader) (r io.ReadCloser, err error) {
	r = lzw.NewReader(r, lzw.MSB, 8)
	return
}

func (c compressorLzw) getCompressionAlgorithm() uint8 {
	return cmpAlgrLzw
}

type compressorZlib struct {
	emptyCompressorOptions
}

func (c compressorZlib) compressor(compressed io.Writer) (w io.WriteCloser, err error) {
	if w, err = zlib.NewWriterLevel(w, zlib.BestCompression); err != nil {
		return
	}
	return
}

func (c compressorZlib) decompressor(compressed io.Reader) (r io.ReadCloser, err error) {
	if r, err = zlib.NewReader(r); err != nil {
		return
	}
	return
}

func (c compressorZlib) getCompressionAlgorithm() uint8 {
	return cmpAlgrZlib
}

type writerNopCloser struct {
	io.Writer
}

func (writerNopCloser) Close() error { return nil }

func newWriterNopCloser(w io.Writer) io.WriteCloser {
	return writerNopCloser{w}
}
