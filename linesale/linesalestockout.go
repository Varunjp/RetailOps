package linesale

import (
	"log"
	"net/http"
	db "retialops/DB"
	"retialops/models"
	responsemodel "retialops/models/responseModel"
	"time"

	"github.com/gin-gonic/gin"
)

// load page

func LineSaleStockOutPage(c *gin.Context){
	c.HTML(http.StatusOK,"linesale_item_stockout.html",nil)
}

// provide items for selected vehicles
func StockOutItems(c *gin.Context){
	var Itemslist []models.LineSale
	var Items []responsemodel.LineSaleStockOut
	vehicle := c.Query("vehicle_id")
	today := time.Now().Format("2006-01-02")
	query := db.Db.Model(&models.LineSale{}).Where("vehicle = ? AND created_at BETWEEN ? AND ? AND status = ?",vehicle,today+" 00:00:00",today+" 23:59:59",false)

	if err := query.Find(&Itemslist).Error; err != nil{
		log.Println("Error while getting data for vehicle :",err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not get items, please try again later"})
		return 
	}

	for _,item := range Itemslist{
		Items = append(Items, responsemodel.LineSaleStockOut{
			ID: item.ID,
			Name: item.ItemName,
			StockIn: item.StockIn,
			StockOut: item.StockOut,
		})
	}

	c.JSON(http.StatusOK,Items)
}

func StockOutUpdate(c *gin.Context){
	var update responsemodel.StockUpdateRequest

	if err := c.ShouldBindJSON(&update); err != nil{
		log.Println("error while binding update line sale items :",err)
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request"})
		return 
	}


	for _, item := range update.StockUpdates {
		var linesale models.LineSale
		if item.StockIn < item.StockOut {
			c.JSON(http.StatusBadRequest,gin.H{"error":"Item out cannot be more than item in"})
			return 
		}

		if err := db.Db.Where("id = ?",item.ItemID).First(&linesale).Error; err != nil{
			log.Println("Failed to get item details :",err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to get item details, please try again later"})
			return 
		}

		linesale.StockOut = item.StockOut
		linesale.Status = true 
		if err := db.Db.Save(&linesale).Error; err !=nil{
			log.Println("Failed to save update :",err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not update item, please try again later"})
			return 
		}
	}

	c.JSON(http.StatusOK,gin.H{"status":true})
}