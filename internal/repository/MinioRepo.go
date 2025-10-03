package repositroy

import "gorm.io/gorm"

type MinioRepo[T any] interface {
	GetData(string) (*gorm.DB, []T)
	SaveData([]T) *gorm.DB
}

type Repo[T any] struct {
	db *gorm.DB
}

func NewRepo[T any](db *gorm.DB) *Repo[T] {
	return &Repo[T]{
		db: db,
	}
}

func (r *Repo[T]) GetData(query string) (*gorm.DB, []T) {
	var result []T
	queryResult := r.db.Where("file_name like ?", query).Take(&result)
	return queryResult, result

}

func (r *Repo[T]) SaveData(dataToSave []T) *gorm.DB {
	saveResult := r.db.Create(dataToSave)
	return saveResult
}
