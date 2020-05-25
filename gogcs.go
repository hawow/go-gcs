package gogcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"io/ioutil"
)

type GoGCSClient interface {
	UploadFile(ctx context.Context, file File) (*UploadedFile, error)
	DownloadFile(ctx context.Context, download DownloadedFile) (*DownloadedFile, error)
}

type GoGSCClient struct {
	Client    *storage.Client
	ProjectID string
	Bucket    string
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
	}
}

func (gcsClient GoGSCClient) UploadFile(ctx context.Context, file File) (*UploadedFile, error) {
	bh := gcsClient.Client.Bucket(gcsClient.Bucket).UserProject(gcsClient.ProjectID)
	obj := bh.Object(GetFullPath(file.Path, file.Name))
	w := obj.NewWriter(ctx)

	defer func() {
		err := gcsClient.Client.Close()
		if err != nil {
			panic(fmt.Errorf("error during closing connection: %v", err))
		}
	}()
	if _, err := io.Copy(w, file.Body); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	if file.IsPublic {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return nil, err
		}
	}

	objAttrs, err := obj.Attrs(ctx)

	if objAttrs == nil {
		return nil, err
	}

	return &UploadedFile{
		Name:        file.Name,
		Size:        objAttrs.Size,
		IsPublic:    file.IsPublic,
		MD5:         MD5BytesToString(objAttrs.MD5),
		Url:         ObjectToUrl(objAttrs),
		ObjectAttrs: objAttrs,
	}, err
}

func (gcsClient GoGSCClient) DownloadFile(ctx context.Context, download DownloadedFile) (*DownloadedFile, error) {
	rc, err := gcsClient.Client.Bucket(gcsClient.Bucket).Object(download.Object).NewReader(ctx)

	defer func() {
		err := gcsClient.Client.Close()
		if err != nil {
			panic(fmt.Errorf("error during closing connection: %v", err))
		}
	}()

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
