package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configurable interface {
	GetConfiguration()
}

type MinioConfig struct {
	Server struct {
		Url  string `yaml:"url"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	Cred struct {
		UserName  string `yaml:"minio-user"`
		UserPass  string `yaml:"minio-pass"`
		AccessKey string `yaml:"access-key"`
		SecretKey string `yaml:"secret-key"`
	} `yaml:"cred"`
}

func (conf *MinioConfig) GetConfiguration() error {

	data, err := os.ReadFile("internal/config/minio-config.yml")
	if err != nil {
		return fmt.Errorf("Error while try to read yaml configuration for minio :: %w", err)
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return fmt.Errorf("Error while try to unmarshal yaml configuration for minio :: %w", err)
	}

	return nil
}
