package client

import (
	"log"
	"testing"
)

func TestGetDrugsComData(t *testing.T) {

	t.Run("Testing status from server ", func(t *testing.T) {
		testCtx := t.Context()
		baseUrl := "https://drugs.com"
		paramPage := "/alpha/a.html"

		client, err := NewDrugsClient(baseUrl)
		if err != nil {
			t.Errorf("error while create new client :: %v", err)
		}
		log.Printf("Success create client :: %+v", client)
		response, err := client.GetData(testCtx, paramPage)
		if err != nil {
			t.Errorf("Error while try test client :: %v", err)
		} else {
			log.Printf("Response result :: %v", response)
		}

	})
}
