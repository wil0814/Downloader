package ftp

import (
	"download/application/model"
	"download/domain/upload/service"
	"download/infrastructure"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"time"
)

type FtpDownload struct {
	fileManager infrastructure.FileManagerInterface
	userFlag    *model.UserFlag
	ftpURL      *FTPUrlInfo
}

func NewFtpDownload(userFlag *model.UserFlag, ftpURL *FTPUrlInfo) *FtpDownload {
	fileManager := infrastructure.NewFileManager()
	return &FtpDownload{
		fileManager: fileManager,
		userFlag:    userFlag,
		ftpURL:      ftpURL,
	}
}

func (d *FtpDownload) Download() error {
	c, err := d.connect()
	if err != nil {
		return fmt.Errorf("%w.\nPlease make sure the FTP URL is in the correct format: ftp://username:password@hostname:port/path/filename", err)
	}

	defer func(c *ftp.ServerConn) {
		err := c.Quit()
		if err != nil {
			log.Println(err)
		}
	}(c)

	err = d.login(c)
	if err != nil {
		return fmt.Errorf("%w.\nPlease make sure the FTP URL is in the correct format: ftp://username:password@hostname:port/path/filename", err)
	}

	err = d.changeDir(c)
	if err != nil {
		return fmt.Errorf("%w.\nPlease make sure the FTP URL is in the correct format: ftp://username:password@hostname:port/path/filename", err)
	}

	return d.retrieveFile(c)
}
func (d *FtpDownload) connect() (*ftp.ServerConn, error) {
	fmt.Println("Connect 嗨嗨")
	addr := d.ftpURL.Host + ":" + d.ftpURL.Port
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	return c, nil
}
func (d *FtpDownload) login(c *ftp.ServerConn) error {
	err := c.Login(d.ftpURL.Username, d.ftpURL.Password)
	if err != nil {
		return err
	}

	return nil
}

func (d *FtpDownload) changeDir(c *ftp.ServerConn) error {
	err := c.ChangeDir(d.ftpURL.Path)
	if err != nil {
		return err
	}

	return nil
}
func (d *FtpDownload) retrieveFile(c *ftp.ServerConn) error {
	body, err := c.Retr(d.ftpURL.Filename)
	if err != nil {
		return fmt.Errorf("%w.\nPlease make sure the FTP URL is in the correct format: ftp://username:password@hostname:port/path/filename", err)
	}
	defer func(body *ftp.Response) {
		err := body.Close()
		if err != nil {
			log.Println(err)
		}
	}(body)

	if string(service.DestinationS3) == d.userFlag.UploadDestination {
		return d.uploadToS3(d.getFileName(), body)
	}

	return d.writeFile(d.getFileName(), body)
}
func (d *FtpDownload) getFileName() string {
	if d.userFlag.FileName == "" {
		return d.ftpURL.Filename
	} else {
		return d.userFlag.FileName
	}
}
func (d *FtpDownload) writeFile(filename string, body io.Reader) error {
	localFile, err := d.fileManager.CreateDestFile(filename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	pb := infrastructure.NewProgressBar()
	bar := pb.CreateBar()

	for {
		buf := make([]byte, 1024)
		n, err := body.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n > 0 {
			bar.SetTotal(bar.Current()+int64(n), false)
			_, err = localFile.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
			bar.IncrBy(n)

		}
		if err == io.EOF {
			break
		}
	}
	// Set the total to -1 to indicate that the download is complete
	bar.SetTotal(-1, true)
	// Wait for all progress bar to complete
	pb.Progress.Wait()

	return nil
}
func (d *FtpDownload) uploadToS3(fileName string, body io.Reader) error {
	u := service.NewUpload(fileName, body, d.userFlag, service.WithS3Uploader())
	return u.Upload(service.DestinationS3)
}
