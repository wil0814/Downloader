package download

import (
	"download/utils"
	"fmt"
)

type Downloader interface {
	Download() error
}

type Download struct {
	FactoryMap map[utils.ProtocolType]DownloaderFactory
}

func NewDownload() *Download {
	d := &Download{
		FactoryMap: make(map[utils.ProtocolType]DownloaderFactory),
	}
	d.registerFactory(utils.ProtocolHTTP, NewHTTPDownloaderFactory())
	d.registerFactory(utils.ProtocolFTP, NewFTPDownloaderFactory())
	return d
}

func (d *Download) registerFactory(protocol utils.ProtocolType, factory DownloaderFactory) {
	d.FactoryMap[protocol] = factory
}
func (d *Download) Download(userFlag utils.UserFlag) error {
	factory, ok := d.FactoryMap[userFlag.Protocol]
	if !ok {
		return fmt.Errorf("unsupported protocol: %s", userFlag.Protocol)
	}

	downloader, err := factory.CreateDownloader(userFlag)
	if err != nil {
		return fmt.Errorf("failed to create downloader: %w", err)
	}
	return downloader.Download()
}
