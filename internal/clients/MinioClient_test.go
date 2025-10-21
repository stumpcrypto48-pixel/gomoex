package client

import (
	"encoding/json"
	"fmt"
	"httpfromtcp/rootmod/internal/config"
	"log"
	"testing"
)

type ObjectForTest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMinioClient(t *testing.T) {
	var cfg config.MinioConfig
	err := cfg.GetConfiguration()
	if err != nil {
		t.Errorf("Error whil try to get minio configuration %v", err)
	}
	minioClient, err := NewMinioClient[[]byte](t.Context(), cfg, "satdata")
	if err != nil {
		t.Errorf("Creating minio client error :: %v", err)
	}
	fileName := "TestFile.json"

	testObject := ObjectForTest{
		Name: "Name",
		Age:  100,
	}

	t.Run("Test minio client save JSON file from io.reader into store", func(t *testing.T) {
		bytesData, err := json.Marshal(testObject)
		if err != nil {
			t.Error("Error while marshal test obj")
		}

		saveInfo, err := minioClient.SaveData(t.Context(), fileName, bytesData)
		if err != nil {
			t.Errorf("error while save data into storage :: %v", err)
		} else {
			log.Printf("Saved obj meta :: %+v", saveInfo)
		}
	})

	t.Run("Test minio client checker", func(t *testing.T) {
		ok, err := minioClient.CheckFileExists(t.Context(), fileName)
		if err != nil {
			t.Errorf("Error while try to check file :: %v", err)
		} else {
			log.Printf("File existing :: %v", ok)
		}
	})

	t.Run("Test minio delete file from bucket", func(t *testing.T) {
		err := minioClient.Delete(t.Context(), fileName)
		if err != nil {
			t.Errorf("Error while try to delete data from S3 :: %v", err)
		} else {
			fmt.Printf("File deleted")
		}
	})
}
