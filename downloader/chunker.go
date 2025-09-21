package downloader

import "fmt"

type Offset struct {
	Start, End uint64
}

func (o Offset) String() string {
	return fmt.Sprintf("bytes=%d-%d", o.Start, o.End)
}

func ChunkN(totalSize uint64, chunkSz uint64) []Offset {
	offsets := []Offset{}

	tempSz := totalSize
	chunkSz = chunkSz - 1

	min := func(a, b uint64) uint64 {
		if a > b {
			return b
		} else {
			return a
		}
	}

	var (
		startIndex uint64 = 0
		endIndex   uint64 = min(chunkSz, totalSize)
	)

	for tempSz >= chunkSz {
		offsets = append(offsets, Offset{Start: startIndex, End: endIndex})
		tempSz -= chunkSz
		startIndex, endIndex = endIndex, endIndex+min(tempSz, chunkSz)
	}

	if totalSize%chunkSz != 0 {
		offsets = append(offsets, Offset{startIndex, endIndex})
	}

	return offsets
}
