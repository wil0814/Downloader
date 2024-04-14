package ftp_download

import (
	"download/download/flag"
	"download/download/ftp_download/ftp_handler"
	"download/file_manager"
	"download/utils/progress_bar"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"time"
)

type FtpDownload struct {
	fileManager file_manager.FileManagerInterface
	flag        *flag.UserFlag
	ftpURL      *ftp_handler.FTPPathInfo
}

func NewFtpDownload(flag *flag.UserFlag, ftpURL *ftp_handler.FTPPathInfo) *FtpDownload {
	fileManager := file_manager.NewFileManager()
	return &FtpDownload{
		fileManager: fileManager,
		flag:        flag,
		ftpURL:      ftpURL,
	}
}

func (d *FtpDownload) Download() error {
	c, err := d.connect()
	if err != nil {
		return fmt.Errorf("%w.\nPlease make sure the FTP URL is in the correct format: ftp://username:password@hostname:port/path/filename", err)
	}

	defer c.Quit()

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
	defer body.Close()
	return d.writeFile(body, d.ftpURL.Filename)
}

func (d *FtpDownload) writeFile(body io.Reader, filename string) error {
	localFile, err := d.fileManager.CreateDestFile(filename)
	if err != nil {
		return err
	}
	defer localFile.Close()

	pb := progress_bar.NewProgressBar()
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
