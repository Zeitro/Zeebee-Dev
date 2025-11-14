package global

import (
	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/aws_s3"
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
		return errors.Wrap(err, "Init S3 Client Error")
	}
	s3Client = cli
	return err
}
