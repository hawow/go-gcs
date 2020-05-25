package gogcs

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"io"
)

type GoGCSClient interface {
	UploadSingleFile(ctx context.Context, file File) (*UploadedFile, error)
}

type GoGSCClient struct {
	Client    *storage.Client
	ProjectID string
	Bucket    string
}

func NewGCSClient(ctx context.Context) *GoGSCClient {
	config := LoadGSCConfig()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.JSONPath))

	if err != nil {
		return nil
	}

	return &GoGSCClient{
		Client:    client,
		ProjectID: config.ProjectID,
		Bucket:    config.Bucket,
	}
}

func (gcsClient GoGSCClient) UploadSingleFile(ctx context.Context, file File) (*UploadedFile, error) {
	bh := gcsClient.Client.Bucket(gcsClient.Bucket).UserProject(gcsClient.ProjectID)
	obj := bh.Object(GetFullPath(file.Path, file.Name))
	w := obj.NewWriter(ctx)

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
