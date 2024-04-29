package s3

import (
	"download/application/model"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

type S3Upload struct {
}

func NewS3Upload() *S3Upload {
	return &S3Upload{}
}

func (s *S3Upload) Upload(fileName string, body io.Reader, userFlag *model.UserFlag) error {
	fmt.Println("Uploading to S3...")
	// 創建新的AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(userFlag.S3Configuration.Region),
		Credentials: credentials.NewStaticCredentials(userFlag.S3Configuration.AccessKey, userFlag.S3Configuration.SecretKey, ""),
	})
	if err != nil {
		return err
	}

	// 創建S3客戶端
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(userFlag.S3Configuration.BucketName),
		Key:    aws.String(fileName),
		Body:   body,
	})

	if err != nil {
		return err
	}

	fmt.Println("Upload complete ^_^")

	return nil
}
