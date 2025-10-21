package main

import (
	"context"
	"httpfromtcp/rootmod/internal/handlers"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {

	mainContext := context.Background()
	server := &http.Server{
		Addr: ":9999",
	}
	log.Printf("Start server at :: %+v", server)

	g, errGroup := errgroup.WithContext(mainContext)

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		<-errGroup.Done()
		return server.Shutdown(mainContext)
	})

	http.HandleFunc("/v1/getSatData", handlers.SatteliteGETHandler)

	if err := g.Wait(); err != nil {
		log.Printf("Exit session :: %v \n", err)
	}
}
