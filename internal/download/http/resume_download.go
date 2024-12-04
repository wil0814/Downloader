package http

//type DownloadParams struct {
//	filename      string
//	contentLength int
//	partSize      int
//	concurrency   int
//}
//
//type ResumeDownload struct {
//	fileManager FileManager
//	timeout     time.Duration
//}
//
//func NewResumeDownload(
//	fileManager FileManager,
//	timeout time.Duration,
//) *ResumeDownload {
//	return &ResumeDownload{
//		fileManager: fileManager,
//		timeout:     timeout,
//	}
//}
//
//func (d *ResumeDownload) Download() error {
//	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
//	defer cancel()
//
//	resp, err := http.Head("http://google.com") // Use HEAD to get content length
//	if err != nil {
//		return fmt.Errorf("failed to fetch headers: %w", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
//	}
//
//	contentLength := int(resp.ContentLength)
//	if contentLength <= 0 {
//		return errors.New("invalid content length")
//	}
//
//	params := DownloadParams{
//		filename:      "123",
//		contentLength: contentLength,
//		partSize:      contentLength / 1,
//	}
//
//	err = d.resumeDownload(ctx, &params)
//	if err != nil {
//		return err
//	}
//
//	// Merge parts
//	if err := d.fileManager.Merge("123", params.concurrency); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (d *ResumeDownload) resumeDownload(ctx context.Context, params *DownloadParams) error {
//	if err := d.fileManager.CreatDir(params.filename); err != nil {
//		return fmt.Errorf("failed to create directory: %w", err)
//	}
//	defer d.cleanupTempFiles(params.filename)
//
//	var wg sync.WaitGroup
//	wg.Add(params.concurrency)
//
//	segmentStart := 0
//	for i := 0; i < params.concurrency; i++ {
//		segmentEnd := segmentStart + params.partSize
//		if i == params.concurrency-1 {
//			segmentEnd = params.contentLength
//		}
//
//		go func(index, start, end int) {
//			defer wg.Done()
//			if err := d.downloadWorker(ctx, params, index, start, end); err != nil {
//				log.Printf("segment %d failed: %v", index, err)
//			}
//		}(i, segmentStart, segmentEnd)
//
//		segmentStart += params.partSize + 1
//	}
//
//	wg.Wait()
//	return nil
//}
//
//func (d *ResumeDownload) downloadWorker(ctx context.Context, params *DownloadParams, index, start, end int) error {
//	req, err := http.NewRequestWithContext(ctx, "GET", "http://google.com", nil)
//	if err != nil {
//		return fmt.Errorf("failed to create request: %w", err)
//	}
//	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return fmt.Errorf("failed to fetch segment %d: %w", index, err)
//	}
//	defer resp.Body.Close()
//
//	partFileName := d.fileManager.GetPartFileName(params.filename, index)
//	partFile, err := d.fileManager.CreateDestFile(partFileName)
//	if err != nil {
//		return fmt.Errorf("failed to create part file: %w", err)
//	}
//	defer partFile.Close()
//
//	_, err = io.Copy(partFile, resp.Body)
//	if err != nil {
//		return fmt.Errorf("failed to write part %d: %w", index, err)
//	}
//
//	return nil
//}
//
//func (d *ResumeDownload) cleanupTempFiles(baseName string) {
//	tempDir := d.fileManager.GetDirName(baseName)
//	if err := os.RemoveAll(tempDir); err != nil {
//		log.Printf("failed to cleanup temp files: %v", err)
//	}
//}
