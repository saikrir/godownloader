package downloader

type Offset struct {
	Start, End uint64
}

func ChunkN(totalSize uint64, chunkSz uint64) []Offset {
	offsets := []Offset{}

	tempSz := totalSize

	var (
		startIndex uint64 = 0
		endIndex   uint64 = chunkSz
	)

	for tempSz >= chunkSz {
		offsets = append(offsets, Offset{Start: startIndex, End: endIndex})
		startIndex, endIndex = endIndex, endIndex+chunkSz
		tempSz -= chunkSz
	}

	return offsets
}
