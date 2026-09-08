package aws_s3

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c = &Config{
	Bucket:       "zeebee-files-staging",
	Region:       "us-west-1",
	AccessKey:    "AKIA4RT6FLQ5RNYTQCPJ",
	AccessSecret: "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",
}

func TestContentTypeMap_Markdown(t *testing.T) {
	assert.Equal(t, "text/markdown", ContentTypeMap["md"])
}

func TestClient_GetPresignUrl(t *testing.T) {
	client, err := NewS3Client(c)
	if err != nil {
		t.Errorf("fail to create s3 client %v", err)
	}

	url, err := client.GetPresignUrl(context.Background(), "test.txt", 120, "txt")

	if err != nil {
		t.Errorf("fail to create presign url:%v", err)
	}
	t.Logf("url:%s", url)
}

func TestClient_UploadFile(t *testing.T) {
	client, err := NewS3Client(c)
	if err != nil {
		t.Errorf("fail to create s3 client %v", err)
	}
	req, err := client.GetPresignUrl(context.Background(), "files-dev/logo.png", 120, "png")
	if err != nil {
		t.Errorf("fail to create presign url:%v", err)
	}
	t.Logf("url:%s", req.URL)
	t.Logf("method:%s", req.Method)
	t.Logf("headers:%v", req.SignedHeader)

	f, err := os.Open("logo.png")
	if err != nil {
		t.Errorf("fail to open file %v", err)
	}
	defer f.Close()

	url, err := client.UploadFile(context.Background(), req, f)

	if err != nil {
		t.Errorf("fail to upload:%v", err)
	}
	t.Logf("s3 file: %s", url)
}
