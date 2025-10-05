package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	"githu.com/saikrir/godownloader/downloader"
	"githu.com/saikrir/godownloader/metadata"
)

// const FileName = "https://i.imgur.com/z4d4kWk.jpg" // Will be a Program Argument eventually
// const FileName = "https://getsamplefiles.com/download/zip/sample-5.zip" // Will be a Program Argument eventually
const FileName = "https://dl.downloadly.ir/Files/Elearning/Vladmihalcea_High-Performance_Java_Persistence_Training_2021-12.part3.rar?nocache=1759618127"

func collectResults(results <-chan downloader.DownloadedChunk, rawBytes []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for result := range results {
		copy(rawBytes[result.Offset.Start:result.Offset.End], result.RawBytes)
	}
}

func main() {
	start := time.Now()
	ctx := context.Background()
	fileMetaData, err := metadata.GetDownloadMetaData(ctx, FileName)

	if err != nil {
		log.Fatalf("failed to get file meta data")
	}
	rawBytes := make([]byte, fileMetaData.Totalsize)
	offsets := downloader.ChunkN(fileMetaData.Totalsize, 8096*10)

	jobQueue := make(chan downloader.DownloadJob)
	results := make(chan downloader.DownloadedChunk)

	nWorkers := runtime.NumCPU()

	wg := new(sync.WaitGroup)
	wg.Add(nWorkers)

	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go collectResults(results, rawBytes, &resultsWg)

	for range nWorkers {
		go downloader.DownloadAsync(ctx, jobQueue, results, wg)
	}
	slog.Info("all workers started")

	for _, offset := range offsets {
		d := downloader.DownloadJob{
			FileURL: FileName,
			Offset:  offset,
		}
		jobQueue <- d
	}
	close(jobQueue)
	wg.Wait()
	close(results)
	resultsWg.Wait()

	fmt.Println("All results done")

	// file, err := downloader.CreateEmptyFile(fileMetaData.Filename, int64(fileMetaData.Totalsize))
	file, err := os.OpenFile(fileMetaData.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("failed to create file %s", err)
	}
	defer file.Close()

	if _, err := file.Write(rawBytes); err != nil {
		slog.Error("failed to write file", "err", err)
	}

	slog.Info("Time Taken for Download ", "Sec", time.Since(start))
}

func seq(ctx context.Context, offsets []downloader.Offset, rawBytes []byte) {
	for _, offset := range offsets {
		cBytes, err := downloader.ChunkDownload(ctx, FileName, offset)
		if err != nil {
			log.Fatalf("failed to download Chunk %s", err)
		}
		slog.Info("File Chunk downloaded ", "chunk", offset, "leng", len(cBytes))
		copy(rawBytes[offset.Start:offset.End], cBytes)
	}
}
