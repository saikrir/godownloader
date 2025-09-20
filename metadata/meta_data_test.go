package metadata

import (
	"context"
	"testing"
)

func TestGetDownload(t *testing.T) {
	input := "https://i.imgur.com/z4d4kWk.jpg"
	metaData, err := GetDownloadMetaData(context.Background(), input)

	if err != nil {
		t.Fatalf("call failed %v \n", err)
	}

	t.Logf("Data gotten %#v", metaData)
}
