package upload

import (
	"download/application/model"
	"io"
)

type Uploader interface {
	Upload(fileName string, body io.Reader, S3Config *model.UserFlag) error
}
