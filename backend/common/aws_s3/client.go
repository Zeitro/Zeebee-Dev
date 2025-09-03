package aws_s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
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

func (c *S3Client) GetPresignUrl(ctx context.Context, objectName string, expire int) (req *v4.PresignedHTTPRequest, err error) {
	presignClient := s3.NewPresignClient(c.client)
	req, err = presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(objectName),
		// ContentType: aws.String("application/octet-stream"),
	}, func(options *s3.PresignOptions) {
		// options.Expires = time.Duration(expire) * time.Second
	})

	return req, nil
}

func (c *S3Client) UploadFile(ctx context.Context, presignedReq *v4.PresignedHTTPRequest, file io.Reader) (url string, err error) {
	var req *http.Request
	var content []byte
	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)
	if err != nil {
		return
	}

	req, err = http.NewRequest(presignedReq.Method, presignedReq.URL, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return
	}

	// 添加必要的请求头
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

	// S3在成功上传时返回200状态码
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应体以提供更多错误信息
		content, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			err = errors.Errorf("http status code: %v", resp.Status)
		} else {
			err = errors.Errorf("http status code: %v, response: %s", resp.Status, string(content))
		}
		return
	}
	content, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Printf("uploand result:%s\n", string(content))

	// 成功上传后返回S3对象的URL
	url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.config.Bucket, c.config.Region, req.URL.Path)
	return url, nil
}
