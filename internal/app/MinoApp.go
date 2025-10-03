package app

import (
	"httpfromtcp/rootmod/internal/services"
)

type MinioApp struct {
	services services.MinioServicer
}

func (app *MinioApp) Start() {
	app.services.Connect()

}
