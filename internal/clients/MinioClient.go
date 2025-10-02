package client

import (
	"context"
	"fmt"
	"httpfromtcp/rootmod/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type CRDClient[T any] interface {
	Create(context.Context, T) error
	Read(context.Context, T) error
	Delete(context.Context, T) error
}

type MinioClient struct {
	client     *minio.Client
	bucketName string
}

func NewMinioClient(c context.Context, cfg config.MinioConfig, bucketName string) (*MinioClient, error) {
	client, err := minio.New(fmt.Sprintf("%v:%v", cfg.Server.Url, cfg.Server.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.Cred.AccessKey, cfg.Cred.SecretKey, ""),
			Secure: false,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Error while creating client :: %v", err)
	}

	return &MinioClient{
		client: client,
	}, nil
}

func (client MinioClient) Create(c context.Context, objectName string) {

}
