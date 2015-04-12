package s3driver

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/koofr/graval"
)

type S3DriverFactory struct {
	Username               string
	Password               string
	AWSRegion              string
	AWSCredentialsProvider aws.CredentialsProvider
	AWSBucketName          string
}

func (f *S3DriverFactory) NewDriver() (d graval.FTPDriver, err error) {
	return &S3Driver{
		Username:               f.Username,
		Password:               f.Password,
		AWSRegion:              f.AWSRegion,
		AWSCredentialsProvider: f.AWSCredentialsProvider,
		AWSBucketName:          f.AWSBucketName,
	}, nil
}
