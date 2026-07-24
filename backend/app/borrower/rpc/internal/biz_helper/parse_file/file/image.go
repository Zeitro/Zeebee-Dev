package file

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// IsImage gates parsing to PNG/JPEG image files.
func IsImage(name string) bool {
	ext := strings.ToLower(name)
	switch {
	case strings.HasSuffix(ext, ".png"), strings.HasSuffix(ext, ".jpg"), strings.HasSuffix(ext, ".jpeg"):
		return true
	}
	return false
}

// ImageContentType maps an image filename to its MIME type for provider upload.
func ImageContentType(name string) string {
	switch {
	case strings.HasSuffix(strings.ToLower(name), ".png"):
		return "image/png"
	default:
		return "image/jpeg"
	}
}

// DownloadImage fetches a presigned image URL and returns its bytes for upload.
func DownloadImage(ctx context.Context, url string) (*[]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("download http %d", res.StatusCode)
	}
	b, err := readCapped(res.Body, maxPDFBytes)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
