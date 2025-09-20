package metadata

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

//Todo initialize logger

type FileMetaData struct {
	Filename     string
	Hash         string
	RangeSupport bool
	Totalsize    uint64
}

func GetDownloadMetaData(ctx context.Context, fileURL string) (FileMetaData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var (
		req          *http.Request
		res          *http.Response
		fileMetaData FileMetaData
		err          error
	)

	if req, err = http.NewRequestWithContext(ctx, http.MethodHead, fileURL, nil); err != nil {
		slog.Error("failed to create request ", "err", err)
		return fileMetaData, err
	}

	req.Header.Set("User-Agent", "curl/8.5.0") // remove later

	if res, err = http.DefaultClient.Do(req); err != nil {
		slog.Error("failed to execute HEAD request ", "url", fileURL, "err", err)
		return fileMetaData, err
	}

	if res.StatusCode != 200 {
		slog.Error("head call did not return successful response code ", "http_status", res.StatusCode)
		return fileMetaData, fmt.Errorf("head call did not return successful response code: %d ", res.StatusCode)
	}

	defer res.Body.Close()

	if fileMetaData, err = collectFileMetaData(res); err != nil {
		return fileMetaData, err
	}

	return fileMetaData, nil
}

func collectFileMetaData(headResponse *http.Response) (FileMetaData, error) {

	fileURL := headResponse.Request.URL.String()
	fileMetaData := FileMetaData{}
	fileMetaData.Filename = path.Base(fileURL)
	fileMetaData.Hash = strings.ReplaceAll(headResponse.Header.Get("etag"), "\"", "")
	fileMetaData.RangeSupport = headResponse.Header.Get("accept-ranges") == "bytes"

	contentLength := headResponse.Header.Get("content-length")
	length, err := strconv.ParseUint(contentLength, 10, 64)

	if err != nil {
		slog.Error("failed to parse content length ", "err", err)
		return FileMetaData{}, err
	}

	slog.Info("MetaData ", "fileName", fileMetaData.Filename, "length", length)

	fileMetaData.Totalsize = length
	return fileMetaData, nil
}
