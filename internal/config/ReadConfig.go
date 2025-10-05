package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigFileReadError struct {
	Message string
	Err     error
}

func (c *ConfigFileReadError) Error() string {
	return fmt.Sprintf("Error while try to read config file :: %v\n", c.Message)
}

func (c *ConfigFileReadError) Unwrap() error {
	return c.Err
}

type Configurable interface {
	GetConfiguration()
}

type PolygonApi struct {
	Creds struct {
		ApiKey string `yaml:"api-key"`
	} `yaml:"creds"`
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

const (
	confFolderPath = "internal/config"
)

func (conf *MinioConfig) GetConfiguration() error {

	data, err := os.ReadFile(fmt.Sprintf("%v/minio-config.yml", confFolderPath))
	if err != nil {
		return fmt.Errorf("Error while try to read yaml configuration for minio :: %w", err)
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return fmt.Errorf("Error while try to unmarshal yaml configuration for minio :: %w", err)
	}

	return nil
}

func (conf *PolygonApi) GetConfiguration() error {
	data, err := os.ReadFile(fmt.Sprintf("%v/api-stock.yml", confFolderPath))
	if err != nil {
		return &ConfigFileReadError{
			Message: err.Error(),
			Err:     err,
		}
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return &ConfigFileReadError{
			Message: err.Error(),
			Err:     err,
		}
	}
	return nil
}
