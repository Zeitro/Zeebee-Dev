package aws_s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (c *S3Client) GetPresignUrl(ctx context.Context, objectName string, expire int) (url string, err error) {
	presignClient := s3.NewPresignClient(c.client)
	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(objectName),
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Duration(expire) * time.Second
	})
	if err != nil {
		return
	}
	return req.URL, nil
}
