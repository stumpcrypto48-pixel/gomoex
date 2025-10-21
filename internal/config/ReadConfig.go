package config

import (
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	confFolderPath = "gomoex/internal/config"
)

func ReadFile(fileName string) ([]byte, error) {
	var relPath string
	currPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Error while try to get current dir path :: %w", err)
	}
	if !strings.Contains(currPath, "internal") {
		relPath = filepath.Join(currPath, "internal", "config", fileName)
	} else {
		relPath = filepath.Join(filepath.Dir(currPath), "config", fileName)
	}

	data, err := os.ReadFile(relPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (conf *MinioConfig) GetConfiguration() error {
	// var relPath string
	// currPath, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("Error while try to get current dir path :: %w", err)
	// }
	// if !strings.Contains(currPath, "internal") {
	// 	relPath = filepath.Join(currPath, "internal", "config", "minio-config.yml")
	// } else {
	// 	relPath = filepath.Join(filepath.Dir(currPath), "config", "minio-config.yml")
	// }

	data, err := ReadFile("minio-config.yml")
	if err != nil {
		dir, errDir := os.Getwd()
		return fmt.Errorf("Error while try to read yaml configuration for minio :: %w | current dir :: %v | err if dir not found :: %v", err, dir, errDir)
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

func NewCerts(root bool) (*x509.CertPool, error) {
	ok := false
	caCertPool := x509.NewCertPool()

	if root {

		rootCert, err := ReadFile("global.pem")
		if err != nil {
			return nil, err
		}
		ok = caCertPool.AppendCertsFromPEM(rootCert)
	} else {
		clientCert, err := ReadFile("dcert.crt")
		if err != nil {
			return nil, err
		}
		ok = caCertPool.AppendCertsFromPEM(clientCert)
	}

	if ok {
		return caCertPool, nil
	}
	return nil, fmt.Errorf("Error while creating CaCertPool")
}
