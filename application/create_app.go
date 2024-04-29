package application

import (
	"download/application/model"
	"download/domain/download/service"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
)

type CreateAPP struct {
	context  *cli.Context
	userFlag *model.UserFlag
}

func NewCreateAPP(context *cli.Context) *CreateAPP {
	return &CreateAPP{context: context}
}
func (c *CreateAPP) Run() error {
	err := c.validateFlags()
	if err != nil {
		return err
	}
	err = c.configureUserFlag()
	if err != nil {
		return err
	}
	err = c.configureS3()
	if err != nil {
		return err
	}
	return c.configureServer()
}
func (c *CreateAPP) validateFlags() error {
	requiredFlags := []string{"url", "ftp"}
	providedCount := 0

	for _, flag := range requiredFlags {
		if c.context.String(flag) != "" {
			providedCount++
		}
	}

	if providedCount == 0 {
		return errors.New("at least one of " + c.joinFlagNames(requiredFlags) + " flag is required")
	} else if providedCount > 1 {
		return errors.New("only one of " + c.joinFlagNames(requiredFlags) + " flag should be provided")
	}

	return nil
}
func (c *CreateAPP) joinFlagNames(flags []string) string {
	names := make([]string, len(flags))
	for i, flag := range flags {
		names[i] = "'" + flag + "'"
	}
	return strings.Join(names, ", ")
}

func (c *CreateAPP) configureUserFlag() error {
	var protocol model.ProtocolType
	fmt.Println("Configuring User Flag")
	c.userFlag = &model.UserFlag{}
	switch {
	case c.context.String("url") != "":
		protocol = model.ProtocolHTTP
		c.userFlag.Path = c.context.String("url")
		break
	case c.context.String("ftp") != "":
		protocol = model.ProtocolFTP
		c.userFlag.Path = c.context.String("ftp")
		break
	default:
		return errors.New("either 'url', 'ftp', or 'sftp' userFlag must be provided")
	}
	c.userFlag.Protocol = protocol
	c.userFlag.FileName = c.context.String("output")
	c.userFlag.Concurrency = c.context.Int("concurrency")
	c.userFlag.UploadDestination = c.context.String("uploadDestination")
	return nil
}
func (c *CreateAPP) configureS3() error {
	if c.userFlag.UploadDestination != "S3" {
		return nil
	}
	fmt.Println("Configuring S3")
	var S3Config model.S3Configuration
	S3Config.BucketName = c.context.String("S3BucketName")
	S3Config.Region = c.context.String("S3Region")
	S3Config.AccessKey = c.context.String("S3AccessKey")
	S3Config.SecretKey = c.context.String("S3SecretKey")

	s3ConfigMap := map[string]*string{
		"S3BucketName": &S3Config.BucketName,
		"S3Region":     &S3Config.Region,
		"S3AccessKey":  &S3Config.AccessKey,
		"S3SecretKey":  &S3Config.SecretKey,
	}
	fieldDescriptions := map[string]string{
		"S3BucketName": "S3 Bucket Name (-b): The name of the S3 bucket to use for storage",
		"S3Region":     "S3 Regions (-r): The AWS region in which the S3 bucket is located",
		"S3AccessKey":  "S3 Access Key (-a): The access key for your AWS account",
		"S3SecretKey":  "S3 Secret Key (-s): The secret key for your AWS account",
	}

	missingFields := []string{}
	for field, value := range s3ConfigMap {
		if *value == "" {
			missingFields = append(missingFields, fieldDescriptions[field])
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing S3 configuration fields: \n%s", strings.Join(missingFields, "\n"))
	}
	c.userFlag.S3Configuration = &S3Config

	return nil
}

func (c *CreateAPP) configureServer() error {
	var d *service.Download

	if c.userFlag.Protocol == model.ProtocolHTTP {
		d = service.NewDownloadService(service.WithHTTPDownloaderFactory(c.userFlag))
		return d.Download(c.userFlag.Protocol)
	}
	if c.userFlag.Protocol == model.ProtocolFTP {
		d = service.NewDownloadService(service.WithFTPDownloader(c.userFlag))
		return d.Download(c.userFlag.Protocol)

	}
	return errors.New("unsupported protocol")
}
