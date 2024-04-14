package http_download

import (
	"download/download/flag"
	"download/file_manager"
	"download/utils/progress_bar"
	"io"
	"net/http"
	"path"
)

type StandardDownload struct {
	userFlag    *flag.UserFlag
	fileManager file_manager.FileManagerInterface
}

func NewStandardDownload(userFlag *flag.UserFlag) *StandardDownload {
	fileManager := file_manager.NewFileManager()
	return &StandardDownload{
		fileManager: fileManager,
		userFlag:    userFlag,
	}
}

func (d *StandardDownload) Download() error {
	return d.retrieveHTTP(d.getFileName())
}

func (d *StandardDownload) getFileName() string {
	if d.userFlag.FileName == "" {
		return path.Base(d.userFlag.Path)
	} else {
		return d.userFlag.FileName
	}
}
func (d *StandardDownload) retrieveHTTP(fileName string) error {
	resp, err := http.Get(d.userFlag.Path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return d.writeResponseToFile(resp, fileName)
}
func (d *StandardDownload) writeResponseToFile(resp *http.Response, filename string) error {
	localFile, err := d.fileManager.CreateDestFile(filename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	pb := progress_bar.NewProgressBar()
	bar := pb.CreateBar()

	for {
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n > 0 {
			bar.SetTotal(bar.Current()+int64(n), false)
			_, err = localFile.Write(buf[:n])
			if err != nil {
				return err
			}
			bar.IncrBy(n)

		}
		if err == io.EOF {
			break
		}
	}

	// Set the total to -1 to indicate that the download is complete
	bar.SetTotal(-1, true)
	pb.Progress.Wait()

	return err
}
