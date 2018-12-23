package films

import (
	_ "fmt"
	"github.com/jinzhu/gorm"
)

type Film struct {
	gorm.Model
	Title        string `gorm:"unique_index"`
	Year         string
	Plot         string `gorm:"size:2048"`
}