package gogcs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type GCSConfig struct {
	Bucket    string `envconfig:"GCS_BUCKET" required:"true"`
	ProjectID string `envconfig:"GCS_PROJECT_ID" required:"true"`
	JSONPath  string `envconfig:"GOOGLE_APPLICATION_CREDENTIALS" required:"true"`
}

func LoadGSCConfig() GCSConfig {
	var config GCSConfig
	//Loads .env file if any
	if err := godotenv.Load(); err != nil {
		//
	}

	envconfig.MustProcess("", &config)
	return config
}
