package routes

import (
	"net/http"
	"retialops/controller"
	"retialops/middleware"

	"github.com/gin-gonic/gin"
)

func GetUrl(router *gin.Engine){
	
	// home
	router.GET("/",middleware.AuthMiddleware(),middleware.NoCaches(),func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,gin.H{"Test":"Successful login"})
	})
	router.GET("/login",middleware.NoCaches(),controller.LoginPage)
	router.POST("/user/login",middleware.NoCaches(),controller.Login)

	//Products
	router.GET("/user/products",middleware.AuthMiddleware(),middleware.NoCaches(),controller.ProductPage)
	router.POST("/user/add-product",middleware.AuthMiddleware(),middleware.NoCaches(),controller.AddProduct)
	router.GET("/product/edit/:id",middleware.AuthMiddleware(),middleware.NoCaches(),controller.EditProductPage)
	router.POST("/user/update-product/:id",middleware.AuthMiddleware(),middleware.NoCaches(),controller.EditProduct)

	// Line sales
	router.GET("/lineSale-items",middleware.AuthMiddleware(),middleware.NoCaches(),controller.LineSalesItemsPage)
	router.GET("/items",controller.GetItems)
	router.POST("/line-sale-items/save",middleware.AuthMiddleware(),middleware.NoCaches(),controller.SaveLineSaleItem)
	router.GET("/line-sale/edit/:id",middleware.AuthMiddleware(),middleware.NoCaches(),controller.EditLineSaleItemPage)
	router.POST("/linesale/update-item/:id",middleware.AuthMiddleware(),middleware.NoCaches(),controller.EditLineSale)
	router.GET("/lineSale-closing",middleware.AuthMiddleware(),middleware.NoCaches(),middleware.CheckLineSaleClosing(),controller.LineSaleClosingPage)

}