package infrastructure

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type FileManagerInterface interface {
	GetPartFileName(filename string, partNum int) string
	GetDirName(fileName string) string
	CreateDestFile(filename string) (*os.File, error)
	CreatDir(filename string) error
	Merge(fileName string, concurrency int) error
}
type FileManager struct {
}

func NewFileManager() *FileManager {
	return &FileManager{}
}
func (d *FileManager) GetPartFileName(filename string, partNum int) string {
	dirName := d.GetDirName(filename)
	return fmt.Sprintf("%s/%s-%d", dirName, filename, partNum)
}
func (d *FileManager) GetDirName(fileName string) string {
	return strings.SplitN(fileName, ".", 2)[0]
}

func (d *FileManager) CreateDestFile(filename string) (*os.File, error) {
	flags := os.O_WRONLY | os.O_CREATE
	destFile, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		return nil, err
	}
	return destFile, nil
}
func (d *FileManager) CreatDir(filename string) error {
	partDir := d.GetDirName(filename)
	err := os.Mkdir(partDir, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (d *FileManager) Merge(fileName string, concurrency int) error {
	flags := os.O_WRONLY | os.O_CREATE

	destFile, err := os.OpenFile(fileName, flags, 0666)

	if err != nil {
		return err
	}
	defer destFile.Close()
	for i := 0; i < concurrency; i++ {
		partFileName := d.GetPartFileName(fileName, i)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(destFile, partFile)
		if err != nil {
			partFile.Close()
			return err
		}
		partFile.Close()

		err = os.Remove(partFileName)
		if err != nil {
			return err
		}
	}
	return nil
}
