package download

import ()

//
//type Download struct {
//	downloadServices map[model.ProtocolType]DownloadInterface
//}
//type Configuration func(d *Download) error
//
//func NewDownloadService(cfg ...Configuration) *Download {
//	d := &Download{
//		downloadServices: make(map[model.ProtocolType]DownloadInterface),
//	}
//
//	for _, c := range cfg {
//		err := c(d)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//	return d
//}
//
//func WithHTTPDownloaderFactory(userFlag *model.UserFlag) Configuration {
//	var service *http.HttpDownload
//	var err error
//	supportResume := checkCanResumeDownload(*userFlag)
//	if supportResume {
//		service, err = http.NewHttpService(http.WithResumeDownload(userFlag))
//	} else {
//		service, err = http.NewHttpService(http.WithStandardDownload(userFlag))
//	}
//	if err != nil {
//		log.Fatal(err)
//	}
//	return withHTTPDownloader(service)
//}
//
//func withHTTPDownloader(d *http.HttpDownload) Configuration {
//	return func(os *Download) error {
//		os.downloadServices[model.ProtocolHTTP] = d.Download
//		return nil
//	}
//}
//
//func WithFTPDownloader(userFlag *model.UserFlag) Configuration {
//	ftpURLParser := ftp2.NewFTPURLParser()
//	ftpURL, err := ftpURLParser.Parse(userFlag.Path)
//	if err != nil {
//		log.Fatal(err)
//	}
//	ftpD := ftp2.NewFtpDownload(userFlag, ftpURL)
//
//	return func(os *Download) error {
//		os.downloadServices[model.ProtocolFTP] = ftpD
//		return nil
//	}
//}
//
//func (d *Download) Download(key model.ProtocolType) error {
//	service, ok := d.downloadServices[key]
//	if !ok {
//		return fmt.Errorf("no download service for protocol: %v", key)
//	}
//
//	err := service.Download()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
