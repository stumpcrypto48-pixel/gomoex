package app

import (
	"fmt"
	"net/http"
)

type PlayerStore struct{}

const jsonContentType = "application/json"

func PlayerServer(response http.ResponseWriter, request *http.Request) {

	fmt.Fprint(response, "20")
}
