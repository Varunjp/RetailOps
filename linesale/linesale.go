package linesale

import (
	"fmt"
	"log"
	"net/http"
	db "retialops/DB"
	"retialops/helper"
	"retialops/models"
	"retialops/models/getitems"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)



func LineSaleStockIn(c *gin.Context) {
	var inventoryReq getitems.InventoryRequest
	tokenStr, _ := c.Cookie("JWT-User")
	empId, _, errTk := helper.DecodeJWT(tokenStr)
	session := sessions.Default(c)

	if errTk != nil{
		log.Println("Error while extracting user details :",errTk)
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Error while fetching user details, please re-logIn."})
		return 
	}

	if err := c.ShouldBindJSON(&inventoryReq); err != nil{
		log.Println("Error while getting inventory data :",err)
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request"})
		return 
	}

	if inventoryReq.Transactions.Vehicle == ""{
		log.Println("Error no vehicle value received")
		c.JSON(http.StatusBadRequest,gin.H{"error":"Vehicle info not passed"})
		return 
	}

	for _,item := range inventoryReq.Items{
		var product models.Product
		db.Db.Preload("ProductDetail").Where("id = ?",item.ItemID).First(&product)

		if product.ProductDetail[0].Stock < int(item.StockIn){
			log.Println("Product not avilable")
			erstr := fmt.Sprintf("%v is not avilable, please update your inventory",product.ItemName)
			c.JSON(http.StatusBadRequest,gin.H{"error":erstr})
			return 
		}

		if item.StockIn < 0 || item.StockIn == 0{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Stock in cannot be less than or equal to zero"})
			return 
		}

	}


	for _,item := range inventoryReq.Items{
		var product models.Product
		db.Db.Preload("ProductDetail").Where("id = ?",item.ItemID).First(&product)
		// need to check if stock need to be reduced while stock in

		lineSale := models.LineSale{
			EmpID: empId,
			Vehicle: inventoryReq.Transactions.Vehicle,
			ProductID: product.ID,
			ItemName: product.ItemName,
			Rate: item.Rate,
			StockIn: int(item.StockIn),
			Status: false,
			Created_at: time.Now(),
		}

		if err := db.Db.Create(&lineSale).Error; err != nil{
			log.Println("Error while saving items :",err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error while saving item, please try again later"})
			return 
		}

	}

	session.Set("message","Items saved successfully")
	session.Save()

	c.JSON(http.StatusOK,nil)
}
