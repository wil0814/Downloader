package upload

import (
	"download/utils"
	"fmt"
	"io"
)

type UploadDestination string

const (
	DestinationS3 UploadDestination = "S3"
	//DestinationFTP UploadDestination = "FTP"
)

type Uploader interface {
	Upload(fileName string, body io.Reader, S3Config *utils.UserFlag) error
}

type Upload struct {
	fileName   string
	body       io.Reader
	userFlag   *utils.UserFlag
	factoryMap map[UploadDestination]UploadFactory
}

func NewUpload(fileName string, body io.Reader, userFlag *utils.UserFlag) *Upload {
	u := &Upload{
		fileName:   fileName,
		body:       body,
		userFlag:   userFlag,
		factoryMap: make(map[UploadDestination]UploadFactory),
	}
	u.registerFactory(DestinationS3, NewS3UploaderFactory())
	//u.registerFactory(DestinationFTP, NewFTPUploaderFactory())
	return u
}
func (u *Upload) registerFactory(destination UploadDestination, factory UploadFactory) {
	u.factoryMap[destination] = factory
}
func (u *Upload) Upload(destination UploadDestination) error {
	factory, ok := u.factoryMap[destination]
	if !ok {
		return fmt.Errorf("unsupported protocol: %s", destination)
	}

	uploader, err := factory.CreateUploader()
	if err != nil {
		return fmt.Errorf("failed to create uploader: %w", err)
	}

	return uploader.Upload(u.fileName, u.body, u.userFlag)
}
