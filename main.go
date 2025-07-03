package main

import (
	"log"
	"os"
	db "retialops/DB"
	"retialops/helper"
	"retialops/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env files")
	}
	db.DbInit()
	helper.CreateSuperUser(db.Db)
	port := os.Getenv("PORT")

	if port == ""{
		port = "3000"
	}

	router := gin.Default()
	router.Use(gin.Logger(),gin.Recovery())

	store := cookie.NewStore([]byte("secret-key"))
	router.Use(sessions.Sessions("Mysession",store))

	// add template helper

	// load static files

	// load html
	router.LoadHTMLGlob("templates/*")

	routes.GetUrl(router)

	router.Run(":"+port)

}