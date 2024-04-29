package http

import (
	"download/application/model"
	"download/domain/download"
)

type Configuration func(d *HttpDownload) error

type HttpDownload struct {
	Download download.DownloadInterface
}

func NewHttpService(cfgs ...Configuration) (*HttpDownload, error) {
	d := &HttpDownload{}
	for _, cfg := range cfgs {
		err := cfg(d)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func WithStandardDownload(userFlag *model.UserFlag) Configuration {
	sd := NewStandardDownload(userFlag)
	return func(d *HttpDownload) error {
		d.Download = sd
		return nil
	}
}

func WithResumeDownload(flag *model.UserFlag) Configuration {
	rd := NewResumeDownload(flag)
	return func(ds *HttpDownload) error {
		ds.Download = rd
		return nil
	}
}
