package facade

import (
	"context"
	client "httpfromtcp/rootmod/internal/clients"
	"httpfromtcp/rootmod/internal/errors"
	"httpfromtcp/rootmod/internal/model"
	"httpfromtcp/rootmod/internal/utils"
	"log"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type SecQueryParams struct {
	Query string `form:"q"`
}

func GetSecDataFacade(c *gin.Context) {

	// Create cancel context
	ctx, cancel := context.WithCancel(c.Request.Context())
	var query SecQueryParams

	// Read request
	if err := c.Bind(&query); err != nil {
		log.Fatalf("Error while binding JSON :: %v", err)
		errors.WriteAPIError(c, err)
		cancel()
		return
	}

	pageNum := atomic.Int64{}
	offset := int64(100)

	workersNuber := 100
	urlChan := make(chan string, workersNuber)
	totalResponse := make([]model.SecModel, 100)

	go func() {
		defer close(urlChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				pageNum.Add(offset)
				urlForSecs, err := client.CreateUrlForSecs(query.Query, pageNum.Load())

				if err != nil {
					log.Printf("Error while creating url for secs :: %v", err)
					cancel()
					errors.WriteAPIError(c, err)
					return
				}
				urlChan <- urlForSecs
			}
		}
	}()

	for {
		newResponseChan := utils.WorkerPool(ctx, urlChan,
			func(urlForSecs string) model.SecModel {
				if response, err := client.GetSecs(ctx, urlForSecs); err != nil {
					log.Printf("Error while getting secs :: %v", err)
					cancel()
					return model.SecModel{}
				} else {
					return response
				}
			})
		for response := range newResponseChan {
			totalResponse = append(totalResponse, response)
		}
		select {
		case <-ctx.Done():
			c.JSON(200, gin.H{"Body": totalResponse})
			return
		default:
		}

	}
}
