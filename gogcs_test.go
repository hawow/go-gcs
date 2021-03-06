package gogcs

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func TestGetFullPath(t *testing.T) {
	t.Run("It should get full path", func(t *testing.T) {
		actualFileName := GetFullPath("/file/path", "file.txt")
		expectedFileName := "/file/path/file.txt"
		if actualFileName != expectedFileName {
			t.Errorf("Expecting result to be %v got %v", expectedFileName, actualFileName)
		}
	})
}

func TestNewGCSClient(t *testing.T) {
	t.Run("It should instantiate new GCS Client", func(t *testing.T) {
		client := NewGCSClient(context.Background())
		if client == nil {
			t.Error("Client is nil")
		}
	})
}

func TestGoGSCClient_UploadFiles(t *testing.T) {
	t.Run("It should upload one file", func(t *testing.T) {
		ctx := context.Background()
		client := NewGCSClient(ctx)
		if client == nil {
			t.Error("Client is nil")
		}
		f, err := os.Open("sample.txt")
		if err != nil {
			log.Fatal(err)
			return
		}
		file := File{
			Name:     "test.txt",
			Path:     "new/test/file",
			Body:     f,
			IsPublic: true,
		}
		var files []File
		files = append(files, file)
		result, err := client.UploadFiles(files)

		if err != nil {
			t.Errorf("error not nil: %v", err)
		}

		if result == nil {
			t.Error("Result should not be nil!")
		}
	})
}

func TestGoGSCClient_DownloadFile(t *testing.T) {
	t.Run("It should download single/Multiple file", func(t *testing.T) {

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()
		client := NewGCSClient(ctx)
		if client == nil {
			t.Error("Client is nil")
			return
		}
		downloadFiles := []DownloadedFile{
			{
				Object: "new/test/file/test.txt",
				Name:   "test.txt",
				Path:   "",
			},
		}
		err := client.DownloadFiles(downloadFiles)

		if err != nil {
			t.Error("error not nil")
		}
	})
}

func TestGoGSCClient_RemoveFile(t *testing.T) {
	t.Run("It should download single/multiple file", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()
		client := NewGCSClient(ctx)
		if client == nil {
			t.Error("Client is nil")
			return
		}
		downloadFiles := []DownloadedFile{
			{
				Object: "new/test/file/test.txt",
				Name:   "test.txt",
				Path:   "",
			},
		}
		err := client.RemoveFiles(downloadFiles)

		if err != nil {
			t.Error("error not nil")
		}
	})
}
