package services

import (
	"context"
	client "httpfromtcp/rootmod/internal/clients"
	"httpfromtcp/rootmod/internal/model"
	"httpfromtcp/rootmod/internal/utils"
	"log"
	"sync/atomic"

	repositroy "httpfromtcp/rootmod/internal/repository"

	"gorm.io/gorm"
)

type SecsService interface {
	GetSecsRequest(c *context.Context, query string) (error, []model.SecModel)
	GetSecsDbQuery(query string) (*gorm.DB, []model.SecModel)
	SaveData([]model.SecModel) *gorm.DB
}

type secsService struct {
	secRepo repositroy.SecsRepo
}

func NewSecService(secRepo repositroy.SecsRepo) *secsService {
	return &secsService{secRepo: secRepo}
}

func (service *secsService) GetSecsRequest(c *context.Context, query string) (error, []model.SecModel) {

	ctx, cancel := context.WithCancel(*c)

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
				urlForSecs, err := client.CreateUrlForSecs(query, pageNum.Load())

				if err != nil {
					log.Printf("Error while creating url for secs :: %v", err)
					cancel()
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
			return nil, totalResponse
		default:
		}
	}
}

func (service *secsService) GetSecsDbQuery(query string) (*gorm.DB, []model.SecModel) {
	return service.secRepo.GetData(query)
}

func (service *secsService) SaveData(saveData []model.SecModel) *gorm.DB {
	return service.secRepo.SaveData(saveData)
}
