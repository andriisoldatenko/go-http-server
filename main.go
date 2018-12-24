package main

import (
	"github.com/andriisoldatenko/go-http-server/auth"
	"github.com/andriisoldatenko/go-http-server/films"
	"github.com/andriisoldatenko/go-http-server/storage"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)


func Migrate(db *gorm.DB) {
	db.AutoMigrate(&films.FilmModel{})
}


func main() {
	db := storage.Init()
	Migrate(db)
	defer db.Close()

	r := gin.Default()
	r.GET("/login", films.LoginHandler)
	r.GET("/implicit/callback", films.AuthCodeCallbackHandler)

	v1 := r.Group("/api")
	v1.Use(auth.AuthMiddleware())
	films.FilmAnonymousRegister(v1.Group("/films"))

	tx1 := db.Begin()
	filmA := films.FilmModel{
		Title: "AAAAAAAAAAAAAAAA",
		Year:  "2018",
		Plot:  "hehddeda",
	}
	tx1.Save(&filmA)
	tx1.Commit()

	r.Run()
}