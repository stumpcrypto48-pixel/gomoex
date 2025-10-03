package services

import (
	"context"
	"errors"
	client "httpfromtcp/rootmod/internal/clients"
	repositroy "httpfromtcp/rootmod/internal/repository"
)

type MinioServicer interface {
	GetData(context.Context, string) error
	PushData(context.Context, string) error
	DeleteData(context.Context, string) error
}

type MinioService[O, T any] struct {
	dbRepository *repositroy.Repo[O]
	minioClient  client.CRDClient[T]
}

func NewMinioService[O, T any](minioRepo *repositroy.Repo[O], client client.CRDClient[T]) *MinioService[O, T] {
	return &MinioService[O, T]{
		minioClient:  client,
		dbRepository: minioRepo,
	}
}

func (service MinioService[O, T]) GetData(ctx context.Context, fileName string) error {
	// 1 check file name in  db
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	_, data := service.dbRepository.GetData(fileName)
	if len(data) == 0 {
		return errors.New("No such file")
	}

	// if exists get data from minio

	service.minioClient.Read(c)
}
