package file_manager

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type FileManager struct {
}

func (d *FileManager) CreatPartDir(filename string) error {
	partDir := d.GetPartDir(filename)
	err := os.Mkdir(partDir, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (d *FileManager) WriteFile(fileName string, partNum int, body io.Reader) error {
	flags := os.O_WRONLY | os.O_CREATE

	partFile, err := os.OpenFile(d.GetPartFileName(fileName, partNum), flags, 0666)
	if err != nil {
		return err
	}
	defer partFile.Close()
	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(partFile, body, buf)

	if err != nil && err != io.EOF {
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

		io.Copy(destFile, partFile)
		partFile.Close()
		os.Remove(partFileName)
	}
	return nil
}

func (d *FileManager) GetPartFileName(filename string, partNum int) string {
	partDir := d.GetPartDir(filename)
	return fmt.Sprintf("%s/%s-%d", partDir, filename, partNum)
}
func (d *FileManager) GetPartDir(fileName string) string {
	return strings.SplitN(fileName, ".", 2)[0]
}
