package main

import (
	"download/download"
	"download/download/flag"
	"errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"strings"
)

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
			var userFlag flag.UserFlag
			var protocol flag.ProtocolType
			switch {
			case c.String("url") != "":
				protocol = flag.ProtocolHTTP
				userFlag.Path = c.String("url")
			case c.String("ftp") != "":
				protocol = flag.ProtocolFTP
				userFlag.Path = c.String("ftp")
			default:
				return errors.New("either 'url', 'ftp', or 'sftp' userFlag must be provided")
			}
			userFlag.Protocol = protocol
			userFlag.FileName = c.String("output")
			userFlag.Concurrency = c.Int("concurrency")
			return d.Download(userFlag)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
