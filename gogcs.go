package gogcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"io/ioutil"
)

type GoGCSClient interface {
	UploadFiles(file []File) ([]UploadedFile, error)
	DownloadFiles(downloads []DownloadedFile) error
	RemoveFiles(downloads []DownloadedFile) error
}

type GoGSCClient struct {
	Client    *storage.Client
	ProjectID string
	Bucket    string
	Context   context.Context
}

func NewGCSClient(ctx context.Context) *GoGSCClient {
	config := LoadGSCConfig()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil
	}

	return &GoGSCClient{
		Client:    client,
		ProjectID: config.ProjectID,
		Bucket:    config.Bucket,
		Context:   ctx,
	}
}

func (gcsClient GoGSCClient) UploadFiles(files []File) ([]UploadedFile, error) {
	bh := gcsClient.Client.Bucket(gcsClient.Bucket).UserProject(gcsClient.ProjectID)

	defer func() {
		err := gcsClient.Client.Close()
		if err != nil {
			panic(fmt.Errorf("error during closing connection: %v", err))
		}
	}()

	var results []UploadedFile

	for _, file := range files {
		obj := bh.Object(GetFullPath(file.Path, file.Name))
		w := obj.NewWriter(gcsClient.Context)

		if _, err := io.Copy(w, file.Body); err != nil {
			return results, err
		}

		if err := w.Close(); err != nil {
			return results, err
		}

		if file.IsPublic {
			if err := obj.ACL().Set(gcsClient.Context, storage.AllUsers, storage.RoleReader); err != nil {
				return results, err
			}
		}

		objAttrs, err := obj.Attrs(gcsClient.Context)

		if objAttrs == nil {
			return results, err
		}

		results = append(results, UploadedFile{
			Name:        file.Name,
			Size:        objAttrs.Size,
			IsPublic:    file.IsPublic,
			MD5:         MD5BytesToString(objAttrs.MD5),
			Url:         ObjectToUrl(objAttrs),
			ObjectAttrs: objAttrs,
		})
	}

	return results, nil

}

func (gcsClient GoGSCClient) downloadFile(download DownloadedFile) (*DownloadedFile, error) {
	rc, err := gcsClient.Client.Bucket(gcsClient.Bucket).Object(download.Object).NewReader(gcsClient.Context)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rc.Close()
		if err != nil {
			panic(fmt.Errorf("error 2 %v", err))
		}
	}()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(GetFullPath(download.Path, download.Name), data, 0644)

	if err != nil {
		return nil, err
	}

	download.Data = &data

	return &download, nil
}

func (gcsClient GoGSCClient) removeFile(download DownloadedFile) error {
	object := gcsClient.Client.Bucket(gcsClient.Bucket).Object(download.Object)

	if err := object.Delete(gcsClient.Context); err != nil {
		return err
	}

	return nil
}

func (gcsClient GoGSCClient) DownloadFiles(downloads []DownloadedFile) error {
	defer func() {
		err := gcsClient.Client.Close()
		if err != nil {
			panic(fmt.Errorf("error during closing connection: %v", err))
		}
	}()
	for k, download := range downloads {
		result, err := gcsClient.downloadFile(download)
		if err != nil {
			return err
		}
		downloads[k].Data = result.Data
	}
	return nil
}

func (gcsClient GoGSCClient) RemoveFiles(downloads []DownloadedFile) error {
	defer func() {
		err := gcsClient.Client.Close()
		if err != nil {
			panic(fmt.Errorf("error during closing connection: %v", err))
		}
	}()

	for _, download := range downloads {
		err := gcsClient.removeFile(download)
		if err != nil {
			return err
		}
	}

	return nil
}
