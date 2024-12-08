package http

import "os"

type FileManager interface {
	CreateDestFile(filename string) (*os.File, error)
	CreatDir(baseName string) error
	GetPartFileName(baseName string, index int) string
	GetDirName(baseName string) string
	Merge(baseName string, concurrency int) error
}

type DownloadService interface {
	Download() error
}
