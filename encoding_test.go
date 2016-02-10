package ictl

import "testing"

func TestEncoding(t *testing.T) {
	lipsum := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	pool := newSlicePool(2000)
	var packet1, packet2, payload *ReusableSlice
	var err error

	if packet1, err = encodeKF(pool, []byte(lipsum), 42); err != nil {
		t.Fatalf("error encoding KF: %v\n", err)
	}

	if _, payload, err = decode(pool, packet1.Slice()); err != nil {
		t.Fatalf("error decoding KF: %v\n", err)
	}

	got := string(payload.Slice())
	if got != lipsum {
		t.Fatalf("compressed then uncompressed data is not equal to original.\norigin: %s\n   got: %s\n", lipsum, got)
	}

	t.Logf("KF compressing/decompressing test passed\n")

	if packet2, err = encodeDF(pool, []byte(lipsum), packet1.Slice(), 42); err != nil {
		t.Fatalf("error encoding DF: %v\n", err)
	}

	if _, payload, err = decode(pool, packet2.Slice()); err != nil {
		t.Fatalf("error decoding DF: %v\n", err)
	}

	data := pool.Get()
	xor(packet1.Slice(), payload.Slice(), data)
	got = string(data.Slice())
	if got != lipsum {
		t.Fatalf("compressed then uncompressed data is not equal to original.\norigin: %s\n   got: %s\n", lipsum, got)
	}

	t.Logf("DF compressing/decompressing test passed\n")
}
