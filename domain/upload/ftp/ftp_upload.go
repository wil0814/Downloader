package ftp

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"time"
)

// TODO 目前還找不到公開的FTP上傳後可檢視 上傳內容來做測試 先暫時棄用
type FTPUpload struct {
}

func NewFTPUpload() *FTPUpload {
	return &FTPUpload{}
}

func (u *FTPUpload) Upload(fileName string, body io.Reader) error {
	fmt.Println("Uploading to FTP Sever...")
	c, err := u.connect()
	if err != nil {
		return err
	}

	defer func(c *ftp.ServerConn) {
		err := c.Quit()
		if err != nil {
			log.Println(err)
		}
	}(c)

	err = u.login(c)
	if err != nil {
		return err
	}

	err = u.changeDir(c)
	if err != nil {
		return err
	}
	//TODO 測試用
	data := bytes.NewBufferString("Hello World")
	err = u.storeFile(c, data)
	if err != nil {

	}
	fmt.Println("Upload complete ^_^")

	return nil
}
func (u *FTPUpload) connect() (*ftp.ServerConn, error) {
	//addr := u.ftpURL.Host + ":" + u.ftpURL.Port
	c, err := ftp.Dial("ftp.speed.hinet.net:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	return c, nil
}
func (u *FTPUpload) login(c *ftp.ServerConn) error {
	err := c.Login("ftp", "ftp")
	if err != nil {
		return err
	}

	return nil
}
func (d *FTPUpload) changeDir(c *ftp.ServerConn) error {
	err := c.ChangeDir("/uploads")
	if err != nil {
		return err
	}

	return nil
}

func (d *FTPUpload) storeFile(c *ftp.ServerConn, body io.Reader) error {
	err := c.Stor("test.txt", body)
	if err != nil {
		return fmt.Errorf("%w.\nPlease make sure the FTP URL is in the correct format: ftp://username:password@hostname:port/path/filename", err)
	}
	return nil
}
