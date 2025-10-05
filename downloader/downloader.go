package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

func CreateEmptyFile(path string, size int64) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	file.Seek(size-1, io.SeekStart)
	file.Write([]byte{0})
	return file, nil
}

type DownloadJob struct {
	fileURL string
	offset  Offset
}

type DownloadedChunk struct {
	rawBytes []byte
	offset   Offset
	err      error
}

func DownloadAsync(ctx context.Context, downloadQueue <-chan DownloadJob, results chan<- DownloadedChunk, wg *sync.WaitGroup) 
	defer wg.Done()
	defer close(result)

	for downloadJob := range downloadQueue {
		rawBytes, err := DownloadChunk(ctx, downloadJob.fileURL, downloadJob.offset)
		downloadChunk := DownloadedChunk{
			rawBytes: rawBytes,
			offset:   downloadJob.offset,
			err:      err,
		}
		result <- downloadChunk
	}
}

func DownloadChunk(ctx context.Context, fileURL string, offset Offset) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)

	var (
		res *http.Response
		err error
	)

	rawBytes := new(bytes.Buffer)

	req.Header.Set("User-Agent", "curl/8.5.0") // remove later
	req.Header.Set("Range", offset.String())

	if res, err = http.DefaultClient.Do(req); err != nil {
		slog.Error("failed to get chunk ", "URL", fileURL, "offset", offset, "err", err)
		return nil, err
	}

	if res.StatusCode != 206 {
		slog.Error("unexpected response code for chunk download ", "fileURL", fileURL, "offset", offset, "statusCode", res.StatusCode)
		return nil, fmt.Errorf("unexpected response code for URL:%s, offset: %s, httpStatus: %d", fileURL, offset, res.StatusCode)
	}

	if _, err = io.Copy(rawBytes, res.Body); err != nil {
		slog.Error("failed to copy raw bytes from response ", "fileURL", fileURL, "offset", offset, "err", err)
		return nil, err
	}

	return rawBytes.Bytes(), nil
}
