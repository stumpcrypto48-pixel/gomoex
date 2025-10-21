package client

import (
	"bytes"
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

type FileOrBytes interface {
	~[]byte
}

type MinioClient[T FileOrBytes] struct {
	client     *minio.Client
	bucketName string
}

func NewMinioClient[T FileOrBytes](c context.Context, cfg config.MinioConfig, bucketName string) (*MinioClient[T], error) {
	client, err := minio.New(fmt.Sprintf("%v:%v", cfg.Server.Url, cfg.Server.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.Cred.AccessKey, cfg.Cred.SecretKey, ""),
			Secure: false,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error while creating client :: %v", err)
	}

	return &MinioClient[T]{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (client MinioClient[T]) CheckFileExists(c context.Context, objectName string) (*minio.ObjectInfo, error) {
	isOnline := client.client.IsOnline()
	if !isOnline {
		return nil, fmt.Errorf("MINIO is not online ")
	}

	objInfo, err := client.client.StatObject(c, client.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &objInfo, nil
}

func (client MinioClient[T]) SaveData(c context.Context, objectName string, data T) (*minio.UploadInfo, error) {
	isOnline := client.client.IsOnline()
	if !isOnline {
		return nil, fmt.Errorf("MINIO is not online ")
	}
	if _, ok := any(data).([]byte); ok {
		reader := bytes.NewReader(data)
		obj, err := client.client.PutObject(c, client.bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{})
		if err != nil {
			return nil, err
		}
		return &obj, nil
	} else {
		return nil, fmt.Errorf("no such type data for upload")
	}
}

func (client MinioClient[T]) Delete(c context.Context, objName string) error {
	err := client.client.RemoveObject(c, client.bucketName, objName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (client MinioClient[T]) Read(c context.Context, objectName string, bucketName string) {

}

func (client MinioClient[T]) Create(c context.Context, objectName string) {

}
