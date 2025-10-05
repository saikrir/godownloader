package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"githu.com/saikrir/godownloader/downloader"
	"githu.com/saikrir/godownloader/metadata"
)

const FileName = "https://i.imgur.com/z4d4kWk.jpg" // Will be a Program Argument eventually

func main() {
	start := time.Now()
	ctx := context.Background()
	fileMetaData, err := metadata.GetDownloadMetaData(ctx, FileName)

	if err != nil {
		log.Fatalf("failed to get file meta data")
	}
	rawBytes := make([]byte, fileMetaData.Totalsize)

	offsets := downloader.ChunkN(fileMetaData.Totalsize, 8096)

	//sequential
	// for _, offset := range offsets {
	// 	cBytes, err := downloader.DownloadChunk(ctx, FileName, offset)
	// 	if err != nil {
	// 		log.Fatalf("failed to download Chunk %s", err)
	// 	}
	// 	slog.Info("File Chunk downloaded ", "chunk", offset, "leng", len(cBytes))
	// 	copy(rawBytes[offset.Start:offset.End], cBytes)
	// }

	// file, err := downloader.CreateEmptyFile(fileMetaData.Filename, int64(fileMetaData.Totalsize))
	// //file, err := os.OpenFile(fileMetaData.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	// if err != nil {
	// 	log.Fatalf("failed to create file %s", err)
	// }
	// defer file.Close()

	// if _, err := file.Write(rawBytes); err != nil {
	// 	slog.Error("failed to write file", "err", err)
	// }

	slog.Info("Time Taken for Download ", "Sec", time.Since(start))
}
