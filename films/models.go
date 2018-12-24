package films

import (
	"github.com/andriisoldatenko/go-http-server/storage"
	"github.com/jinzhu/gorm"
)

type FilmModel struct {
	gorm.Model
	Title        string `gorm:"unique_index"`
	Year         string
	Plot         string `gorm:"size:2048"`
}

func FindOneFilm() (FilmModel, error) {
	db := storage.GetDB()
	var model FilmModel
	tx := db.Begin()
	tx.First(&model, 1)
	err := tx.Commit().Error
	return model, err
}