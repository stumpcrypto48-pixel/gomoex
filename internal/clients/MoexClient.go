package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"httpfromtcp/rootmod/internal/model"
	"log"
	"net/http"
	"net/url"
)

func CreateUrlForSecs(query string, pageNum int64) (string, error) {
	baseUrl := "https://iss.moex.com/iss/securities.xml"
	u, err := url.Parse(baseUrl)
	if err != nil {
		log.Printf("Error while parsing an url")
		return "", err
	}

	q := u.Query()

	q.Set("q", query)
	q.Set("start", fmt.Sprintf("%d", pageNum))

	u.RawQuery = q.Encode()

	finalURL := u.String()

	return finalURL, nil
}

func GetSecs(c context.Context, urlForSecs string) (model.SecModel, error) {

	var emptyModel model.SecModel

	log.Printf("Making GET request to :: %s", urlForSecs)
	resp, err := http.Get(urlForSecs)
	if err != nil {
		log.Printf("Error while getting response from url :: %v", err)
		return emptyModel, err
	}
	defer resp.Body.Close()

	dat := xml.NewDecoder(resp.Body)
	var response model.SecModel
	if err := dat.Decode(&response); err != nil {
		log.Printf("Error while decoding xml :: %v", err)
		return emptyModel, err
	}

	if len(response.Data.Rows.Items) == 0 {
		log.Printf("End of reading")
		return emptyModel, fmt.Errorf("Custom error end of requests")
	}

	return response, nil
}
