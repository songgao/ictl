package ictl

import "testing"

func TestEncoding(t *testing.T) {
	lipsum := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	lipsum2 := lipsum[:42] + "." + lipsum[43:]

	pool := newSlicePool(2000)
	var packet1, packet2, payload1, payload2 *ReusableSlice
	var err error

	if packet1, err = encode(pool, []byte(lipsum), 42, frameKF, CAAuto); err != nil {
		t.Fatalf("error encoding KF: %v\n", err)
	}

	t.Logf("KF encoded: %d bytes (from %d bytes); header: %x\n", len(packet1.Slice()), len(lipsum), packet1.Slice()[:4])

	if _, payload1, err = decode(pool, packet1.Slice()); err != nil {
		t.Fatalf("error decoding KF: %v\n", err)
	}
	packet1.Done()

	got := string(payload1.Slice())
	if got != lipsum {
		t.Fatalf("compressed then uncompressed data is not equal to original.\norigin: %s\n   got: %s\n", lipsum, got)
	}

	t.Logf("KF compressing/decompressing test passed\n")

	payload2 = pool.Get()
	xor(payload1.Slice(), []byte(lipsum2), payload2)
	if packet2, err = encode(pool, payload2.Slice(), 42, frameDF, CAAuto); err != nil {
		t.Fatalf("error encoding DF: %v\n", err)
	}
	payload2.Done()

	t.Logf("DF encoded: %d bytes (from %d bytes); header: %x\n", len(packet2.Slice()), len(lipsum2), packet2.Slice()[:4])

	if _, payload2, err = decode(pool, packet2.Slice()); err != nil {
		t.Fatalf("error decoding DF: %v\n", err)
	}
	packet2.Done()

	data := pool.Get()
	xor(payload1.Slice(), payload2.Slice(), data)
	payload1.Done()
	payload2.Done()
	got = string(data.Slice())
	if got != lipsum2 {
		t.Fatalf("compressed then uncompressed data is not equal to original.\norigin: %s\n   got: %s\n", lipsum, got)
	}

	t.Logf("DF compressing/decompressing test passed\n")
}
