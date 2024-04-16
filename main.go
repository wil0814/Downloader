package main

import (
	"download/download"
	"download/utils"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"strings"
)

func main() {
	concurrency := runtime.NumCPU()
	d := download.NewDownload()

	app := &cli.App{
		Name:  "download",
		Usage: "download files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Usage:   "url to download",
			},
			&cli.StringFlag{
				Name:    "ftp",
				Aliases: []string{"f"},
				Usage:   "ftp to download",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output file",
			},
			&cli.StringFlag{
				Name:    "uploadDestination",
				Aliases: []string{"d"},
				Usage:   "upload destination",
			},
			&cli.StringFlag{
				Name:    "S3BucketName",
				Aliases: []string{"b"},
				Usage:   "S3 Bucket Name",
			},
			&cli.StringFlag{
				Name:    "S3Region",
				Aliases: []string{"r"},
				Usage:   "S3 Regions",
			},
			&cli.StringFlag{
				Name:    "S3AccessKey",
				Aliases: []string{"a"},
				Usage:   "S3 Access Key",
			},
			&cli.StringFlag{
				Name:    "S3SecretKey",
				Aliases: []string{"s"},
				Usage:   "S3 Secret Key",
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"n"},
				Value:   concurrency,
				Usage:   "concurrency number",
			},
		},
		Action: func(c *cli.Context) error {
			err := validateFlags(c)
			if err != nil {
				return err
			}
			var userFlag utils.UserFlag
			var protocol utils.ProtocolType
			switch {
			case c.String("url") != "":
				protocol = utils.ProtocolHTTP
				userFlag.Path = c.String("url")
			case c.String("ftp") != "":
				protocol = utils.ProtocolFTP
				userFlag.Path = c.String("ftp")
			default:
				return errors.New("either 'url', 'ftp', or 'sftp' userFlag must be provided")
			}
			userFlag.Protocol = protocol
			userFlag.FileName = c.String("output")
			userFlag.Concurrency = c.Int("concurrency")
			userFlag.UploadDestination = c.String("uploadDestination")

			err = configureS3(c, userFlag)
			if err != nil {
				return err
			}

			return d.Download(userFlag)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
func validateFlags(c *cli.Context) error {
	requiredFlags := []string{"url", "ftp"}
	providedCount := 0

	for _, flag := range requiredFlags {
		if c.String(flag) != "" {
			providedCount++
		}
	}

	if providedCount == 0 {
		return errors.New("at least one of " + joinFlagNames(requiredFlags) + " flag is required")
	} else if providedCount > 1 {
		return errors.New("only one of " + joinFlagNames(requiredFlags) + " flag should be provided")
	}

	return nil
}
func joinFlagNames(flags []string) string {
	names := make([]string, len(flags))
	for i, flag := range flags {
		names[i] = "'" + flag + "'"
	}
	return strings.Join(names, ", ")
}
func configureS3(c *cli.Context, userFlag utils.UserFlag) error {
	if userFlag.UploadDestination != "S3" {
		return nil
	}
	var S3Config utils.S3Configuration
	S3Config.BucketName = c.String("S3BucketName")
	S3Config.Region = c.String("S3Region")
	S3Config.AccessKey = c.String("S3AccessKey")
	S3Config.SecretKey = c.String("S3SecretKey")

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
	userFlag.S3Configuration = &S3Config

	return nil
}
