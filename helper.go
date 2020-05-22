package gogcs

import (
	"cloud.google.com/go/storage"
	"fmt"
)

func ObjectToUrl(objAttrs *storage.ObjectAttrs) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", objAttrs.Bucket, objAttrs.Name)
}

func MD5BytesToString(bytes []byte) string {
	return fmt.Sprintf("%x", bytes)
}

func GetFullPath(path string, name string) string {
	return fmt.Sprintf("%s/%s", path, name)
}
