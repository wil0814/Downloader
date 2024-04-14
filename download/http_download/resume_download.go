package http_download

import (
	"download/download/flag"
	"download/download/http_download/http_handler"
	"download/file_manager"
	"download/utils/progress_bar"
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
	userFlag    *flag.UserFlag
}

func NewResumeDownload(client http_handler.HttpHandlerInterface, file file_manager.FileManagerInterface, flag *flag.UserFlag) *ResumeDownload {
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

	defer resp.Body.Close()
	return d.resumeDownload(fileName, int(resp.ContentLength))
}
func (d *ResumeDownload) resumeDownload(filename string, contentLength int) error {
	concurrency := d.userFlag.Concurrency
	partSize := contentLength / concurrency
	err := d.fileManager.CreatDir(filename)
	if err != nil {
		return err
	}
	defer os.RemoveAll(d.fileManager.GetDirName(filename))

	pb := progress_bar.NewProgressBar()
	bar := pb.CreateBar()
	params := &DownloadParams{
		filename:      filename,
		contentLength: contentLength,
		partSize:      partSize,
		bar:           bar,
	}
	d.initiateDownloadWorkers(params)

	// Set the total to -1 to indicate that the download is complete
	bar.SetTotal(-1, true)
	// Wait for all progress bar to complete
	pb.Progress.Wait()

	err = d.fileManager.Merge(filename, concurrency)
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
	defer resp.Body.Close()

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
