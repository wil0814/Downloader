package download

import (
	"download/download/flag"
	"fmt"
)

type Downloader interface {
	Download() error
}

type Download struct {
	FactoryMap map[flag.ProtocolType]DownloaderFactory
}

func NewDownload() *Download {
	d := &Download{
		FactoryMap: make(map[flag.ProtocolType]DownloaderFactory),
	}
	d.registerFactory(flag.ProtocolHTTP, NewHTTPDownloaderFactory())
	d.registerFactory(flag.ProtocolFTP, NewFTPDownloaderFactory())
	return d
}

func (d *Download) registerFactory(protocol flag.ProtocolType, factory DownloaderFactory) {
	d.FactoryMap[protocol] = factory
}
func (d *Download) Download(flag flag.UserFlag) error {
	factory, ok := d.FactoryMap[flag.Protocol]
	if !ok {
		return fmt.Errorf("unsupported protocol: %s", flag.Protocol)
	}

	downloader, err := factory.CreateDownloader(flag)
	if err != nil {
		return fmt.Errorf("failed to create downloader: %w", err)
	}

	return downloader.Download()
}
