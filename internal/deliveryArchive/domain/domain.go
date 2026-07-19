package delivery

import "errors"

var (
	ErrNoFile     = errors.New("no file to deliver")
	ErrSendFailed = errors.New("failed to send file")
)

type FileItem struct {
	Path string
	Name string
}
