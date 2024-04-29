package _interface

import (
	"download/application"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

type CLI struct {
}

func NewCLI() *CLI {
	return &CLI{}
}
func (c *CLI) Run() error {
	concurrency := runtime.NumCPU()

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
			app := application.NewCreateAPP(c)
			return app.Run()
		},
	}
	return app.Run(os.Args)

}
