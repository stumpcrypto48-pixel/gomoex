package main

import (
	"context"
	"fmt"
	"httpfromtcp/rootmod/internal/config"
	"log"

	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// get minio configuration
	// ctx := context.Background()
	var cfg config.MinioConfig
	err := cfg.GetConfiguration()
	if err != nil {
		log.Fatalf("Error while startup app :: %v", err)
	}

	log.Printf("Minio configuration :: %+v", cfg)

	client, err := minio.New(fmt.Sprintf("%v:%v", cfg.Server.Url, cfg.Server.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.Cred.AccessKey, cfg.Cred.SecretKey, ""),
			Secure: false,
		},
	)
	if err != nil {
		log.Fatalf("Error while try to connect to minio :: %w", err)
	}
	log.Printf("Client :: %+v", client)

	ctx := context.Background()

	bucketList, err := client.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("Error while try to get bucket list from minio :: %v", err)
	}

	// log.Printf("Available buckets in minio :: %+v", bucketList)
	videoBucket, err := findBucketByName("videos", bucketList)
	if err != nil {
		log.Fatalf("Error while try to get bucket by name :: %v", err)
	}
	log.Printf("Video bucket :: %+v", videoBucket)

	info, err := client.FPutObject(ctx,
		videoBucket.Name,
		"video_example.avi",
		"Files/file_example_MP4_1920_18MG.mp4",
		minio.PutObjectOptions{
			ContentType: "video/mp4",
		})
	if err != nil {
		log.Fatalf("Error while try to save video into s3 :: %v", err)
	}

	log.Printf("Video uploaded :: %+v", info)

}

func findBucketByName(bucketName string, bucketList []minio.BucketInfo) (minio.BucketInfo, error) {
	var result minio.BucketInfo

	for _, bucket := range bucketList {
		if bucket.Name == bucketName {
			return bucket, nil
		}
	}
	return result, fmt.Errorf("No bucket with name %v presented", bucketName)

}
