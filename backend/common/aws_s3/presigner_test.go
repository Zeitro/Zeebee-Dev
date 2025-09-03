package aws_s3

import (
	"context"
	"os"
	"testing"
)

func TestClient_GetPresignUrl(t *testing.T) {
	c := &Config{
		Bucket:       "zeebee-files-staging",
		Region:       "us-west-1",
		AccessKey:    "AKIAZYJFTPUSWNDYUDER",
		AccessSecret: "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",
	}
	client, err := NewS3Client(c)
	if err != nil {
		t.Errorf("failed to create s3 client %v", err)
	}

	url, err := client.GetPresignUrl(context.Background(), "test.txt", 120)

	if err != nil {
		t.Errorf("failed to create presign url:%v", err)
	}
	t.Logf("url:%s", url)
}

func TestClient_UploadFile(t *testing.T) {
	c := &Config{
		Bucket:       "zeebee-files-staging",
		Region:       "us-west-1",
		AccessKey:    "AKIAZYJFTPUSWNDYUDER",
		AccessSecret: "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",
	}
	client, err := NewS3Client(c)
	if err != nil {
		t.Errorf("failed to create s3 client %v", err)
	}
	req, err := client.GetPresignUrl(context.Background(), "files-dev/test.png", 120)
	if err != nil {
		t.Errorf("failed to create presign url:%v", err)
	}
	t.Logf("url:%s", req.URL)
	t.Logf("method:%s", req.Method)
	t.Logf("headers:%v", req.SignedHeader)

	f, err := os.Open("test.png")
	if err != nil {
		t.Errorf("failed to open file %v", err)
	}
	defer f.Close()

	url, err := client.UploadFile(context.Background(), req, f)

	if err != nil {
		t.Errorf("failed to upload:%v", err)
	}
	t.Logf("%s", url)
}
