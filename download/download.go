package download

import (
	"download/file_manager"
	"download/http_client"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

type Download struct {
	client      *http_client.HTTPClient
	file        *file_manager.FileManager
	concurrency int
}

func NewDownload(concurrency int) *Download {
	return &Download{
		client:      &http_client.HTTPClient{},
		file:        &file_manager.FileManager{},
		concurrency: concurrency,
	}
}

func (d *Download) Download(url string, filename string) error {
	if filename == "" {
		filename = path.Base(url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Header.Get("Accept-Ranges"))

	if d.client.IsHttpStatusOK(resp) && d.client.SupportRange(resp) {
		//fmt.Println("content length", resp.ContentLength)
		return d.multiDownload(url, filename, int(resp.ContentLength))
	}

	return d.singleDownload(resp, filename)
}

// TODO
func (d *Download) singleDownload(resp *http.Response, filename string) error {
	flags := os.O_WRONLY | os.O_CREATE
	destFile, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"singleDownload downloading",
	)

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(destFile, bar), resp.Body, buf)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// TODO  multi progress bar
func (d *Download) multiDownload(url string, filename string, contentLength int) error {
	partSize := contentLength / d.concurrency

	err := d.file.CreatPartDir(filename)
	if err != nil {
		return err
	}
	defer os.RemoveAll(d.file.GetPartDir(filename))

	d.startDownloadHandlers(url, filename, contentLength, partSize)

	err = d.file.Merge(filename, d.concurrency)
	if err != nil {
		return err
	}

	return nil
}

func (d *Download) startDownloadHandlers(url string, filename string, contentLength int, partSize int) {
	var wg sync.WaitGroup
	wg.Add(d.concurrency)
	//fmt.Println("concurrency", d.concurrency)

	start := 0
	for partNumber := 0; partNumber < d.concurrency; partNumber++ {
		go d.downloadHandler(url, filename, contentLength, partSize, &wg, partNumber, start)
		start += partSize + 1
	}
	wg.Wait()
}

func (d *Download) downloadHandler(url string, filename string, contentLength int, partSize int, wg *sync.WaitGroup, partNumber int, start int) {
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
}
