package main

import (
	"httpfromtcp/rootmod/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	log.Println("Starting server...")
	router.GET("/getSecData", handlers.GetSecData)
	router.Run(":8081")
}
