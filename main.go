package main

import (
	"download/download"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
)

func main() {
	concurrency := runtime.NumCPU()

	app := &cli.App{
		Name:  "download",
		Usage: "download files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Aliases:  []string{"u"},
				Usage:    "url to download",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "output file",
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"n"},
				Value:   concurrency,
				Usage:   "concurrency number",
			},
		},
		Action: func(c *cli.Context) error {
			strURL := c.String("url")
			filename := c.String("output")
			concurrency := c.Int("concurrency")
			d, err := download.NewDownload(strURL, concurrency)
			if err != nil {
				return err
			}
			return d.Download(strURL, filename)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
