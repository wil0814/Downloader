package upload

import (
	"download/upload/s3"
)

type UploadFactory interface {
	CreateUploader() (Uploader, error)
}

type S3UploaderFactory struct{}

func NewS3UploaderFactory() *S3UploaderFactory {
	return &S3UploaderFactory{}
}
func (f *S3UploaderFactory) CreateUploader() (Uploader, error) {
	return s3.NewS3Upload(), nil
}

//type FTPUploaderFactory struct{}
//
//func NewFTPUploaderFactory() *FTPUploaderFactory {
//	return &FTPUploaderFactory{}
//}
//func (f *FTPUploaderFactory) CreateUploader() (Uploader, error) {
//	return ftp.NewFTPUpload(), nil
//}
