package routes

import (
	"net/http"
	"retialops/controller"
	"retialops/linesale"
	"retialops/middleware"
	"retialops/vehicles"

	"github.com/gin-gonic/gin"
)

func GetUrl(router *gin.Engine){
	
	// home
	router.GET("/",middleware.AuthMiddleware(),middleware.NoCaches(),func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK,gin.H{"Test":"Successful login"})
	})
	router.GET("/login",middleware.NoCaches(),controller.LoginPage)
	router.POST("/user/login",middleware.NoCaches(),controller.Login)
	router.GET("/logout",middleware.NoCaches(),controller.Logout)

	//Products
	router.GET("/user/products",middleware.AuthMiddleware(),middleware.NoCaches(),controller.ProductPage)
	router.POST("/user/add-product",middleware.AuthMiddleware(),middleware.NoCaches(),controller.AddProduct)
	router.GET("/product/edit/:id",middleware.AuthMiddleware(),middleware.NoCaches(),controller.EditProductPage)
	router.POST("/user/update-product/:id",middleware.AuthMiddleware(),middleware.NoCaches(),controller.EditProduct)

	// Line sales
	router.GET("/lineSale-items",middleware.AuthMiddleware(),middleware.NoCaches(),controller.LineSalesItemsPage)
	router.GET("/items",controller.GetItems)
	router.POST("/line-sale-items/save",middleware.AuthMiddleware(),middleware.NoCaches(),controller.SaveLineSaleItem)
	router.GET("/lineSale-stockout",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.LineSaleStockOutPage)
	router.GET("/smrd/vehicles-list",middleware.AuthMiddleware(),middleware.NoCaches(),vehicles.VehicleList)
	router.GET("/smrd/linesale-stockout",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.StockOutItems)
	router.POST("/smrd/linesale-stockout",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.StockOutUpdate)
	router.GET("/line-sale/edit/:id",middleware.AuthMiddleware(),middleware.NoCaches(),controller.EditLineSaleItemPage)
	router.POST("/linesale/update-item/:id",middleware.AuthMiddleware(),middleware.NoCaches(),controller.EditLineSale)
	router.GET("/lineSale-vypar",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.VyparPage)
	router.GET("/smrd/vyapar-stockout",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.VyaparStockOut)
	router.POST("/smrd/linesale-vyapar",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.VyaparUpdate)
	router.GET("lineSale-sales",middleware.AuthMiddleware(),middleware.NoCaches(),middleware.CheckLineSaleClosing(),linesale.LinesaleSalesPage)
	router.GET("/linesale/totalamount",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.LineSaleTotal)
	router.POST("/linesale/sales",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.LineSaleSubmit)
	router.GET("/lineSale-expense",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.ExpensePage)
	router.GET("/smrd/linesale-cash",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.ExpensePerVehicle)
	router.POST("/smrd/submit-expenses",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.ExpenseItems)
	router.GET("/lineSale-closing",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.LineSaleClosingPage)
	router.GET("/smrd/line-sale-closing",middleware.AuthMiddleware(),middleware.NoCaches(),linesale.LineSaleClosingDetails)
	
	// Vehicles
	router.GET("/smrd/vehicles",middleware.AuthMiddleware(),middleware.AuthSuperUser(),middleware.NoCaches(),vehicles.VehiclesPage)
	router.POST("/smrd/add-vehicle",middleware.AuthMiddleware(),middleware.AuthSuperUser(),middleware.NoCaches(),vehicles.AddVehicle)
}