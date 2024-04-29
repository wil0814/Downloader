package model

type ProtocolType string

const (
	ProtocolHTTP ProtocolType = "http"
	ProtocolFTP  ProtocolType = "ftp"
)

type UserFlag struct {
	Path              string
	Protocol          ProtocolType
	FileName          string
	Concurrency       int
	UploadDestination string
	S3Configuration   *S3Configuration
}
type S3Configuration struct {
	Region     string
	BucketName string
	AccessKey  string
	SecretKey  string
}
