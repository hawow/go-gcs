package gogcs

import (
	"os"
	"testing"
)

func TestLoadGSCConfig(t *testing.T) {
	config := LoadGSCConfig()

	t.Run("Load Config File from ENV", func(t *testing.T) {

		if config.ProjectID != os.Getenv("GCS_PROJECT_ID") {
			t.Errorf("Expecting Project ID to be %v, got %v", os.Getenv("GCS_PROJECT_ID"), config.ProjectID)
		}

		if config.Bucket != os.Getenv("GCS_BUCKET") {
			t.Errorf("Expecting Bucket to be %v, got %v", os.Getenv("GCS_BUCKET"), config.Bucket)
		}
	})
}
