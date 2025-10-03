package repositroy

import (
	"httpfromtcp/rootmod/internal/model"

	"gorm.io/gorm"
)

type SecsRepo interface {
	GetData(query string) (*gorm.DB, []model.Row)
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

func (repo *secRepo) GetData(query string) (*gorm.DB, []model.Row) {
	var result []model.Row
	secModelResult := repo.db.Where("name like ? ", query).Take(&result)
	return secModelResult, result
}

func (repo *secRepo) SaveData(resultToSave []model.SecModel) *gorm.DB {
	rows := make([]model.Row, len(resultToSave))
	for _, secModel := range resultToSave {
		rows = append(rows, secModel.Data.Rows.Items...)
	}
	saveResult := repo.db.Create(rows)
	return saveResult
}
