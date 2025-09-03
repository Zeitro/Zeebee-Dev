package aws_s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var accessKey = "AKIAZYJFTPUSWNDYUDER"
var accessSecret = "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"

func GetPresignUrl(ctx context.Context, expire int64) (url string, err error) {
	// 使用静态凭证创建配置
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, accessSecret, "")),
	)
	if err != nil {
		return "", err
	}

	s3Client := s3.NewFromConfig(cfg)

	presignClient := s3.NewPresignClient(s3Client)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:  aws.String("zeebee-files-staging"),
		Key:     aws.String("/files-dev"),
		Expires: aws.Time(time.Now().Add(time.Duration(expire) * time.Second)),
	})
	if err != nil {
		return
	}
	return req.URL, err
}
