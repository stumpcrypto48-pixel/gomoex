package facade

import (
	"context"
	"httpfromtcp/rootmod/internal/errors"
	"httpfromtcp/rootmod/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SecQueryParams struct {
	Query string `form:"q"`
}

type GetSecDataFacade interface {
	GetData(c *gin.Context)
}

type secDataFacade struct {
	service services.SecsService
}

func NewSecDataFacade(service services.SecsService) *secDataFacade {
	return &secDataFacade{service: service}
}

func (facade *secDataFacade) GetData(c *gin.Context) {

	// Parse input json
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var inputJson SecQueryParams
	if err := c.Bind(&inputJson); err != nil {
		log.Printf("Error while try to parse input param for %v", c.Request)
		errors.WriteAPIError(c, err)
		return
	}
	// Check DB for DATA

	db, result := facade.service.GetSecsDbQuery(inputJson.Query)
	if db.Error != nil {
		log.Printf("Error while try to get data from database %v", db.Error)
		errors.WriteAPIError(c, db.Error)
		return
	}

	// If data exists get data from database
	if len(result) != 0 {
		c.JSON(http.StatusOK, gin.H{"response": result})
		return
	}
	// if not load from MOEX API
	err, moexResult := facade.service.GetSecsRequest(&ctx, inputJson.Query)
	if err != nil {
		errors.WriteAPIError(c, err)
	}
	// and save it into database (async)
	errSaveChan := make(chan error)
	go func() {
		db := facade.service.SaveData(moexResult)
		if db.Error != nil {
			log.Printf("Error while try to save data into database :: %v", db.Error)
			errSaveChan <- db.Error
			return
		}
	}()

	// and return to user
	c.JSON(http.StatusOK, gin.H{"response": moexResult})
}
