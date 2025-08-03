package linesale

import (
	"log"
	"net/http"
	db "retialops/DB"
	"retialops/models"
	"retialops/models/getitems"
	"time"

	"github.com/gin-gonic/gin"
)

func ExpensePage(c *gin.Context){
	var Vehicle []models.Vehicle

	if err := db.Db.Find(&Vehicle).Error; err != nil{
		log.Println("Could not load vehicle details :",err)
		c.HTML(http.StatusInternalServerError,"linesale_expense.html",gin.H{"error":"Falied to fetch vehicle details,please try again later"})
		return 
	}

	c.HTML(http.StatusOK,"linesale_expense.html",gin.H{"Vehicle":Vehicle})

}

func ExpensePerVehicle(c *gin.Context){
	var Sale models.LineSaleClosing
	vID := c.Query("vehicle_id")
	today := time.Now().Format("2006-01-02")

	if err := db.Db.Where("status = ? AND vehicle = ? AND created_at BETWEEN ? AND ? ",false,vID,today+" 00:00:00",today+" 23:59:59").First(&Sale).Error; err != nil{
		log.Println("Failed to fetch amount :",err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch today sale details,please try again later"})
		return 
	}


	if Sale.Cash < 1 {
		c.JSON(http.StatusBadRequest,gin.H{"error":"No sales listed for the selected vehicle"})
		return 
	}

	c.JSON(http.StatusOK,gin.H{"cash":Sale.Cash})
}

func ExpenseItems(c *gin.Context){
	var expReq getitems.ExpenseRequest
	var Sale models.LineSaleClosing

	if err := c.ShouldBindJSON(&expReq); err != nil{
		log.Println("Error while getting expense items :",err)
		c.JSON(http.StatusBadRequest,gin.H{"error":"Failed to get item, please try again later"})
		return 
	}

	today := time.Now().Format("2006-01-02")

	if err := db.Db.Where("vehicle = ? AND created_at BETWEEN ? AND ? ",expReq.VehicleID,today+" 00:00:00",today+" 23:59:59").First(&Sale).Error; err != nil{
		log.Println("Failed to fetch amount :",err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to fetch today sale details,please try again later"})
		return 
	}

	if Sale.Status{
		log.Println("Sale already closed for today")
		c.JSON(http.StatusBadRequest,gin.H{"error":"Sale already closed for today, cannot add expenses"})
		return 
	}

	for _,item := range expReq.Transaction{

		expItem := models.LineSaleExpenses{
			LineSaleID: Sale.ID,
			Type: item.Item,
			Amount: item.Amount,
			Created_at: time.Now(),
		}

		if err := db.Db.Create(&expItem).Error; err != nil{
			log.Println("Failed to save expense item :",err)
			DeleteAllExpenseToday(Sale.ID)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to save expense, please try again later"})
			return 
		}
	}

	Sale.Status = true

	if err := db.Db.Save(&Sale).Error; err != nil{
		log.Println("Failed to save sales update :",err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to save update status in DB"})
		return 
	}

	c.JSON(http.StatusOK,gin.H{"message":"Expense items saved successfully","status":Sale.Status})
}


func DeleteAllExpenseToday(lineID uint ){
	today := time.Now().Format("2006-01-02")
	var TodayExpense []models.LineSaleExpenses

	if err := db.Db.Where("line_sale_id = ? AND created_at BETWEEN ? AND ? ",lineID,today+" 00:00:00",today+" 23:59:59").Find(&TodayExpense).Error; err != nil{
		log.Println("Could not load todays data :",err)
	}

	for _,item := range TodayExpense{
		if err := db.Db.Delete(&item).Error; err != nil{
			log.Println("Could not delete today's data :",err)
		}
	}
}