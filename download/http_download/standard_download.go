package http_download

import (
	"download/file_manager"
	"download/upload"
	"download/utils"
	"io"
	"log"
	"net/http"
	"path"
)

type StandardDownload struct {
	userFlag    *utils.UserFlag
	fileManager file_manager.FileManagerInterface
}

func NewStandardDownload(userFlag *utils.UserFlag) *StandardDownload {
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if string(upload.DestinationS3) == d.userFlag.UploadDestination {
		return d.uploadToS3(fileName, resp.Body)
	}

	return d.writeResponseToFile(resp, fileName)
}
func (d *StandardDownload) writeResponseToFile(resp *http.Response, filename string) error {
	localFile, err := d.fileManager.CreateDestFile(filename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	pb := utils.NewProgressBar()
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
func (d *StandardDownload) uploadToS3(fileName string, body io.Reader) error {
	u := upload.NewUpload(fileName, body, d.userFlag)
	return u.Upload(upload.DestinationS3)
}
