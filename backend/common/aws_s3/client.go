package aws_s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"zeebee-backend/common/errorx"
)

type S3Client struct {
	client *s3.Client
	config *Config
}

type Config struct {
	Bucket       string
	Region       string
	AccessKey    string
	AccessSecret string
}

var ContentTypeMap = map[string]string{
	"jpeg": "image/jpeg",
	"jpg":  "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"pdf":  "application/pdf",
	"doc":  "application/msword",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
}

func NewS3Client(conf *Config) (client *S3Client, err error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(conf.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.AccessSecret, "")),
	)
	if err != nil {
		return
	}

	return &S3Client{
		client: s3.NewFromConfig(cfg),
		config: conf,
	}, nil
}

func (c *S3Client) GetPresignUrl(ctx context.Context, objectName string, expire int, suffix string) (req *v4.PresignedHTTPRequest, err error) {
	suffix = strings.ToLower(suffix)
	contentType, ok := ContentTypeMap[suffix]
	if !ok {
		err = errorx.Errorf("Invalid File Type:%s", suffix)
		return nil, err
	}

	presignClient := s3.NewPresignClient(c.client)
	req, err = presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(objectName),
		// ContentType: aws.String("application/octet-stream"),
		ContentType: aws.String(contentType),
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Duration(expire) * time.Second
	})

	return req, nil
}

func (c *S3Client) UploadFile(ctx context.Context, presignedReq *v4.PresignedHTTPRequest, file io.Reader) (url string, err error) {
	var req *http.Request
	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)
	if err != nil {
		return
	}

	req, err = http.NewRequestWithContext(ctx, presignedReq.Method, presignedReq.URL, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return
	}

	// set header
	for s, strings := range presignedReq.SignedHeader {
		for _, v := range strings {
			req.Header.Add(s, v)
		}

	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	req.Header.Set("Content-Length", strconv.Itoa(buf.Len()))
	defer resp.Body.Close()

	// failed to upload
	if resp.StatusCode != http.StatusOK {
		// get error response
		content, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			err = errorx.Errorf("HTTP Status Code: %v", resp.Status)
		} else {
			err = errorx.Errorf("HTTP Status Code: %v, Response: %s", resp.Status, string(content))
		}
		return
	}

	// return url
	url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.config.Bucket, c.config.Region, strings.TrimPrefix(req.URL.Path, "/"))
	return url, nil
}
