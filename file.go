package gogcs

import (
	"cloud.google.com/go/storage"
	"io"
)

type File struct {
	Path     string
	Name     string
	Body     io.Reader
	IsPublic bool
}

type UploadedFile struct {
	Name        string
	MD5         string
	IsPublic    bool
	Url         string
	Size        int64
	ObjectAttrs *storage.ObjectAttrs
}

type DownloadedFile struct {
	Object string
	Name   string
	Path   string
	Data   *[]byte
}
