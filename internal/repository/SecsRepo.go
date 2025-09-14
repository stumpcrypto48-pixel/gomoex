package repositroy

import (
	"httpfromtcp/rootmod/internal/model"

	"gorm.io/gorm"
)

type SecsRepo interface {
	GetData(query string) (*gorm.DB, []model.SecModel)
	SaveData(resultToSave []model.SecModel) *gorm.DB
}

type secRepo struct {
	db *gorm.DB
}

func NewSecRepo(db *gorm.DB) *secRepo {
	return &secRepo{
		db: db,
	}
}

func (repo *secRepo) GetData(query string) (*gorm.DB, []model.SecModel) {
	var result []model.SecModel
	secModelResult := repo.db.Where("name like \"%\" + ? + \"%\"", query).Take(&result)
	return secModelResult, result
}

func (repo *secRepo) SaveData(resultToSave []model.SecModel) *gorm.DB {
	saveResult := repo.db.Create(resultToSave)
	return saveResult
}
