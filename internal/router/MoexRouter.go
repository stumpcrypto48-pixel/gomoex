package router

import (
	"httpfromtcp/rootmod/internal/facade"
	repository "httpfromtcp/rootmod/internal/repository"
	repositroy "httpfromtcp/rootmod/internal/repository"
	"httpfromtcp/rootmod/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MoexRouter(db *gorm.DB, moexRouter *gin.RouterGroup) {
	var (
		secModelRepo repository.SecsRepo     = repositroy.NewSecRepo(db)
		secService   services.SecsService    = services.NewSecService(secModelRepo)
		secFacade    facade.GetSecDataFacade = facade.NewSecDataFacade(secService)
	)

	moexRouter.GET("/getSecData", secFacade.GetData)

}
