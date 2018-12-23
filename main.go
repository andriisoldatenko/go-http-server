package main

import (
	"github.com/andriisoldatenko/go-http-server/storage"
	"github.com/gin-gonic/gin"

	"github.com/andriisoldatenko/go-http-server/films"
	"github.com/jinzhu/gorm"
)


func Migrate(db *gorm.DB) {
	db.AutoMigrate(&films.Film{})
}


func main() {
	db := storage.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()

	//v1 := r.Group("/api")
	//v1.Use()
	r.Run()
}