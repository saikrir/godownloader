package downloader

import "testing"

func TestChunkN0(t *testing.T) {
	totalSz := 1025
	chunkSz := 1000

	chunks := ChunkN(uint64(totalSz), uint64(chunkSz))

	t.Logf("Chunks %+v", chunks)

	if len(chunks) != 2 {
		t.Fatalf("expected 2, got %d ", len(chunks))
	}
}

func TestChunkN(t *testing.T) {
	totalSz := 146515
	chunkSz := 8096

	chunks := ChunkN(uint64(totalSz), uint64(chunkSz))

	for _, chnk := range chunks {
		t.Logf("Offset %s, length: %d \n", chnk, chnk.End-chnk.Start)
	}

	if len(chunks) != 19 {
		t.Fatalf("expected 2, got %d ", len(chunks))
	}
}

func TestChunkN2(t *testing.T) {
	totalSz := 999
	chunkSz := 1000

	chunks := ChunkN(uint64(totalSz), uint64(chunkSz))

	t.Logf("Chunks %+v", chunks)

	if len(chunks) != 1 {
		t.Fatalf("expected 2, got %d ", len(chunks))
	}
}
func TestChunkN3(t *testing.T) {
	totalSz := 1000
	chunkSz := 1000

	chunks := ChunkN(uint64(totalSz), uint64(chunkSz))

	t.Logf("Chunks %+v", chunks)

	if len(chunks) != 2 {
		t.Fatalf("expected 2, got %d ", len(chunks))
	}
}
