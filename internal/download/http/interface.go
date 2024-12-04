package http

import "os"

type FileManager interface {
	CreateDestFile(filename string) (*os.File, error)
	GetPartFileName(baseName string, index int) string
	CreatDir(baseName string) error
	Merge(baseName string, concurrency int) error
	GetDirName(baseName string) string
}

type DownloadService interface {
	Download() error
}
