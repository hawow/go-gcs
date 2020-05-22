package gogcs

import (
	"context"
	"log"
	"os"
	"testing"
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

func TestGoGSCClient_UploadSingleFile(t *testing.T) {
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
		result, err := client.UploadSingleFile(ctx, file)

		if err != nil {
			t.Errorf("error not nil: %v", err)
		}

		if result == nil {
			t.Error("Result should not be nil!")
		}
	})
}
