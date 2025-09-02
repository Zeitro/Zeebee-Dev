package aws_s3

import (
	"context"
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

	url, err := client.GetPresignUrl(context.Background(), "test.txt")

	if err != nil {
		t.Errorf("failed to create presign url:%v", err)
	}
	t.Logf("url:%s", url)
}
