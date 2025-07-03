package routes

import (
	"retialops/controller"

	"github.com/gin-gonic/gin"
)

func GetUrl(router *gin.Engine){
	router.GET("/",controller.LoginPage)
}