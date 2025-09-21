package downloader

import "testing"

func TestChunkN(t *testing.T) {
	totalSz := 1025
	chunkSz := 1000

	chunks := ChunkN(uint64(totalSz), uint64(chunkSz))

	t.Logf("Chunks %#v", chunks)

	if len(chunks) != 2 {
		t.Fatalf("expected 2, got %d ", len(chunks))
	}

}
