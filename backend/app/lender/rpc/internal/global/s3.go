package global

import (
	"zeebee-backend/app/lender/config"
	"zeebee-backend/common/aws_s3"
	"zeebee-backend/common/errorx"
)

var s3Client *aws_s3.S3Client

func GetS3Client() *aws_s3.S3Client {
	return s3Client
}

func initS3Client(cnf config.AwsS3Conf) (err error) {

	cli, err := aws_s3.NewS3Client(&aws_s3.Config{
		Bucket:       cnf.Bucket,
		Region:       cnf.Region,
		AccessKey:    cnf.AccessKey,
		AccessSecret: cnf.AccessSecret,
	})

	if err != nil {
		return errorx.Wrap(err, "Init S3 Client Error")
	}
	s3Client = cli
	return err
}
