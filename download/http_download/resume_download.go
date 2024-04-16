package http_download

import (
	"download/download/http_download/http_handler"
	"download/file_manager"
	"download/upload"
	"download/utils"
	"github.com/vbauerster/mpb/v8"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

type DownloadParams struct {
	filename      string
	contentLength int
	partSize      int
	bar           *mpb.Bar
}
type ResumeDownload struct {
	client      http_handler.HttpHandlerInterface
	fileManager file_manager.FileManagerInterface
	userFlag    *utils.UserFlag
}

func NewResumeDownload(client http_handler.HttpHandlerInterface, file file_manager.FileManagerInterface, flag *utils.UserFlag) *ResumeDownload {
	return &ResumeDownload{
		client:      client,
		fileManager: file,
		userFlag:    flag,
	}
}

func (d *ResumeDownload) Download() error {
	return d.retrieveHTTP(d.getFileName())
}
func (d *ResumeDownload) getFileName() string {
	if d.userFlag.FileName == "" {
		return path.Base(d.userFlag.Path)
	} else {
		return d.userFlag.FileName
	}
}
func (d *ResumeDownload) retrieveHTTP(fileName string) error {
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
	return d.resumeDownload(fileName, int(resp.ContentLength))
}

func (d *ResumeDownload) resumeDownload(fileName string, contentLength int) error {
	concurrency := d.userFlag.Concurrency
	partSize := contentLength / concurrency
	err := d.fileManager.CreatDir(fileName)
	if err != nil {
		return err
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Println(err)
		}
	}(d.fileManager.GetDirName(fileName))

	pb := utils.NewProgressBar()
	bar := pb.CreateBar()
	params := &DownloadParams{
		filename:      fileName,
		contentLength: contentLength,
		partSize:      partSize,
		bar:           bar,
	}
	d.initiateDownloadWorkers(params)

	// Set the total to -1 to indicate that the download is complete
	bar.SetTotal(-1, true)
	// Wait for all progress bar to complete
	pb.Progress.Wait()

	err = d.fileManager.Merge(fileName, concurrency)
	if err != nil {
		return err
	}

	if string(upload.DestinationS3) == d.userFlag.UploadDestination {
		return d.uploadToS3(fileName)
	}

	if err != nil {
		return err
	}
	return nil
}

func (d *ResumeDownload) initiateDownloadWorkers(params *DownloadParams) {
	concurrency := d.userFlag.Concurrency
	var wg sync.WaitGroup
	wg.Add(concurrency)

	segmentStart := 0
	for segmentIndex := 0; segmentIndex < concurrency; segmentIndex++ {

		segmentEnd := segmentStart + params.partSize
		if segmentIndex == concurrency-1 {
			segmentEnd = params.contentLength
		}

		go d.downloadWorker(params, &wg, segmentIndex, segmentStart, segmentEnd)

		segmentStart += params.partSize + 1
	}
	wg.Wait()
}

func (d *ResumeDownload) downloadWorker(params *DownloadParams, wg *sync.WaitGroup, segmentIndex int, segmentStart int, segmentEnd int) {
	defer wg.Done()

	resp, err := d.client.DownloadPartial(d.userFlag.Path, segmentStart, segmentEnd)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	partFile, err := d.fileManager.CreateDestFile(d.fileManager.GetPartFileName(params.filename, segmentIndex))
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()

	for {
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {

		}

		if n > 0 {
			params.bar.SetTotal(params.bar.Current()+int64(n), false)
			_, err := partFile.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
			params.bar.IncrBy(n)
		}

		if err == io.EOF {
			break
		}
	}
}
func (d *ResumeDownload) uploadToS3(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println(err)
		}
	}(fileName)
	defer file.Close()
	u := upload.NewUpload(fileName, file, d.userFlag)
	return u.Upload(upload.DestinationS3)
}
