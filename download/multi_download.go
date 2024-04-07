package download

import (
	"download/file_manager"
	"download/http_client"
	"download/utils/progress_bar"
	"github.com/vbauerster/mpb/v8"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

type MultiDownloader struct {
	client      *http_client.HTTPClient
	file        *file_manager.FileManager
	concurrency int
}

func NewMultiDownloader(client *http_client.HTTPClient, file *file_manager.FileManager, concurrency int) *MultiDownloader {
	return &MultiDownloader{
		client:      client,
		file:        file,
		concurrency: concurrency,
	}
}
func (d *MultiDownloader) Download(url string, filename string) error {
	if filename == "" {
		filename = path.Base(url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	return d.multiDownload(url, filename, int(resp.ContentLength))
}
func (d *MultiDownloader) multiDownload(url string, filename string, contentLength int) error {
	partSize := contentLength / d.concurrency
	err := d.file.CreatDir(filename)
	if err != nil {
		return err
	}
	defer os.RemoveAll(d.file.GetDirName(filename))

	pb := progress_bar.NewProgressBar()
	bar := pb.CreateBar(int64(contentLength))

	d.startDownloadHandlers(url, filename, contentLength, partSize, bar)

	err = d.file.Merge(filename, d.concurrency)
	if err != nil {
		return err
	}

	return nil
}

func (d *MultiDownloader) startDownloadHandlers(url string, filename string, contentLength int, partSize int, bar *mpb.Bar) {
	var wg sync.WaitGroup
	wg.Add(d.concurrency)

	start := 0
	for partNumber := 0; partNumber < d.concurrency; partNumber++ {
		go d.downloadHandler(url, filename, contentLength, partSize, &wg, partNumber, start, bar)
		start += partSize + 1
	}
	wg.Wait()
}

func (d *MultiDownloader) downloadHandler(url string, filename string, contentLength int, partSize int, wg *sync.WaitGroup, partNumber int, start int, bar *mpb.Bar) {
	defer wg.Done()
	end := start + partSize
	if partNumber == d.concurrency-1 {
		end = contentLength
	}

	resp, err := d.client.DownloadPartial(url, start, end)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	err = d.file.WriteFile(filename, partNumber, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bar.IncrInt64(int64(partSize))
}
