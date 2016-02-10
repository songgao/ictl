package ictl

import (
	"bytes"
	"io"
	"testing"
)

func TestCompressors(t *testing.T) {
	lipsum := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	var err error

	for _, v := range compressors {
		cmp := v()

		compressed := new(bytes.Buffer)
		var w io.WriteCloser
		if w, err = cmp.compressor(compressed); err != nil {
			t.Fatalf("error creating compressor (%d): %v\n", cmp.getCompressionAlgorithm(), err)
		}
		if _, err = w.Write([]byte(lipsum)); err != nil {
			t.Fatalf("error writing into the compressor (%d): %v\n", cmp.getCompressionAlgorithm(), err)
		}
		if err = w.Close(); err != nil {
			t.Fatalf("error closing the compressor (%d): %v\n", cmp.getCompressionAlgorithm(), err)
		}

		uncompressed := new(bytes.Buffer)
		var r io.ReadCloser
		if r, err = cmp.decompressor(compressed); err != nil {
			t.Fatalf("error creating decompressor (%d): %v\n", cmp.getCompressionAlgorithm(), err)
		}
		if _, err = uncompressed.ReadFrom(r); err != nil {
			t.Fatalf("error reading from the compressor (%d): %v\n", cmp.getCompressionAlgorithm(), err)
		}
		if err = r.Close(); err != nil {
			t.Fatalf("error closing the compressor (%d): %v\n", cmp.getCompressionAlgorithm(), err)
		}

		got := string(uncompressed.Bytes())
		if got != lipsum {
			t.Fatalf("compressed then uncompressed data is not equal to original. origin: %s; got: %s\n", lipsum, got)
		}

		t.Logf("compressor (%d) test passed\n", cmp.getCompressionAlgorithm())
	}
}
