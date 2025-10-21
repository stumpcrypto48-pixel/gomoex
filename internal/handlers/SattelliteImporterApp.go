package handlers

import (
	"context"
	"fmt"
	client "httpfromtcp/rootmod/internal/clients"
	"httpfromtcp/rootmod/internal/config"
	"io"
	"log"
	"net/http"
	"time"
)

func SatteliteGETHandler(response http.ResponseWriter, request *http.Request) {
	// response.Header().Set("Content-Type", "text/html; charset=ascii")
	// response.Header().Set("Access-Control-Allow-Origin", "*")
	// response.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	// fmt.Fprint(response, "10")

	log.Printf("Start query for sat data")

	rCtx := request.Context()

	var minioCfg config.MinioConfig
	err := minioCfg.GetConfiguration()
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	bucketName := "satdata"
	fileName := fmt.Sprintf("sattelite_data_%v.json", time.Now().Truncate(24*time.Hour).Format("2006-01-02"))

	fileExists, err := CheckFileExists(rCtx, bucketName, fileName, minioCfg)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	if fileExists {
		log.Printf("File already exists")
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("File already exists no need new file creation"))
		return
	}

	httpClient, err := http.NewRequestWithContext(rCtx, http.MethodGet, "https://celestrak.org/NORAD/elements/gp.php?GROUP=active&FORMAT=json", nil)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	clientResp, err := http.DefaultClient.Do(httpClient)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	defer clientResp.Body.Close()

	body, err := io.ReadAll(clientResp.Body)
	if err != nil {
		fmt.Fprint(response, err)
		return
	}

	err = SatSaveDataIntoStore(rCtx, body, bucketName, fileName)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		log.Printf("error while try to save into store :: %v", err.Error())
		return
	}

	response.WriteHeader(http.StatusOK)
	fmt.Fprintf(response, "Data stored success")
	// var ommResponse []sgp4.OMM

	// _ = json.Unmarshal(body, &ommResponse)
	// // if err != nil {
	// // 	fmt.Fprint(response, err)
	// // 	return
	// // }

	// var satNames []model.SatNames

	// _ = json.Unmarshal(body, &satNames)
	// fmt.Println(satNames)

	// resultLLA := services.OMMConverter(ommResponse, satNames)

	// LLAjson, err := json.Marshal(resultLLA)
	// if err != nil {
	// 	fmt.Fprint(response, err)
	// }

	// response.WriteHeader(http.StatusOK)
	// response.Write(LLAjson)
}

func CheckFileExists(c context.Context, bucketName string, fileName string, mCfg config.MinioConfig) (bool, error) {

	ctx, cancel := context.WithTimeout(c, time.Second*30)
	defer cancel()

	s3client, err := client.NewMinioClient[[]byte](ctx, mCfg, bucketName)
	if err != nil {
		return false, err
	}
	_, err = s3client.CheckFileExists(ctx, fileName)
	if err != nil {
		return false, nil
	} else {
		return true, nil
	}
}

func SatSaveDataIntoStore(c context.Context, satData []byte, bucketName, fileName string) error {

	ctx, cancel := context.WithTimeout(c, time.Second*30)
	defer cancel()

	var cfg config.MinioConfig
	err := cfg.GetConfiguration()
	if err != nil {
		return err
	}

	s3client, err := client.NewMinioClient[[]byte](ctx, cfg, bucketName)
	if err != nil {
		return err
	}

	meta, err := s3client.SaveData(ctx, fileName, satData)
	if err != nil {
		return err
	}
	fmt.Printf("Data save success meta :: %+v", meta)

	return nil
}
