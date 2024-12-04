package upload

//type UploadDestination string
//
//const (
//	DestinationS3  UploadDestination = "S3"
//	DestinationFTP UploadDestination = "FTP"
//)
//
//type Upload struct {
//	fileName       string
//	body           io.Reader
//	userFlag       *model.UserFlag
//	uploadServices map[UploadDestination]Uploader
//}
//type Configuration func(d *Upload) error
//
//func NewUpload(fileName string, body io.Reader, userFlag *model.UserFlag, cfgs ...Configuration) *Upload {
//	u := &Upload{
//		fileName:       fileName,
//		body:           body,
//		userFlag:       userFlag,
//		uploadServices: make(map[UploadDestination]Uploader),
//	}
//	for _, cfg := range cfgs {
//		err := cfg(u)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//	return u
//}
//
//func WithS3Uploader() Configuration {
//	return func(os *Upload) error {
//		os.uploadServices[DestinationS3] = NewS3Upload()
//		return nil
//	}
//}
//
//func (u *Upload) Upload(destination UploadDestination) error {
//	uploader, ok := u.uploadServices[destination]
//	if !ok {
//		return fmt.Errorf("unsupported protocol: %s", destination)
//	}
//
//	return uploader.Upload(u.fileName, u.body, u.userFlag)
//}
