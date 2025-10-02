package services

import (
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type MinioServicer interface {
	GetConfig()
	GetData()
	PushData()
	DeleteData()
}

type MinioService struct {
	dbRepository *gorm.DB
	minioClient  minio.Client
}
