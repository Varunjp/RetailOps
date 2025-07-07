package controller

import (
	"log"
	"math"
	"net/http"
	db "retialops/DB"
	"retialops/helper"
	"retialops/models"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type inventoryItem struct{
	ItemID			string 		`json:"itemId"`
	ItemName		string 		`json:"itemName"`
	Rate 			float64		`json:"rate"`
	StockIn			float64		`json:"stockIn"`
	StockOut 		float64		`json:"stockOut"`
	SaleItem 		float64		`json:"saleItem"`
	Damage 			float64		`json:"damage"`
}

type inventoryRequest struct{
	Items		[]inventoryItem 	`json:"items"`
}

func LineSalesItemsPage(c *gin.Context){

	session := sessions.Default(c)
	lineError := session.Get("lineError")
	pageStr := c.DefaultQuery("page","1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1{
		page = 1
	}

	today := time.Now().Format("2006-01-02")

	const pageSize = 10
	var total int64
	db.Db.Model(&models.LineSale{}).Where("created_at BETWEEN ? AND ? ",today+" 00:00:00",today+" 23:59:59").Count(&total)

	totalPages := int(math.Ceil(float64(total)/float64(pageSize)))

	if page > totalPages && totalPages != 0 {
		page = totalPages
	}
	offset := (page - 1)* pageSize
	var lineSales []models.LineSale
	if err := db.Db.Where("created_at >= ? AND created_at <= ? ",today+" 00:00:00",today+" 23:59:59").Limit(pageSize).Offset(offset).Find(&lineSales).Error; err != nil{
		if err != gorm.ErrRecordNotFound{
			log.Println("Failed to retrive data from db :",err)
			c.HTML(http.StatusInternalServerError,"linesales.html",gin.H{"error":"Failed to retrive today's data"})
			return 
		}
	}

	type response struct {
		ID 			uint 
		Name 		string
		Rate 		float64
		Stock 		int 
	}

	var lineSalesResponse []response
	for _,item := range lineSales{
		lineSalesResponse = append(lineSalesResponse, response{
			ID: item.ID,
			Name: item.ItemName,
			Rate: item.Rate,
			Stock: item.StockOut,
		})
	}

	var pages []int 
	startPage := page - 2
	if startPage < 1{
		startPage = 1
	}
	endPage := startPage+4
	if endPage > totalPages{
		endPage = totalPages
	}

	for i := startPage; i<= endPage;i++{
		pages = append(pages, i)
	}

	if lineError != nil{
		session.Delete("lineError")
		session.Save()
	}

	c.HTML(http.StatusOK,"linesales.html",gin.H{
		"products":lineSalesResponse,
		"CurrentPage":page,
		"TotalPages":totalPages,
		"Pages":pages,
		"error":lineError,
	})
}

func GetItems(c *gin.Context){
	var Products []models.Product

	if err := db.Db.Preload("ProductDetail").Where("status = ?","Active").Find(&Products).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"Failed to get products"})
		return
	}

	type Item struct{
		ID 		uint	`json:"id"`
		Name 	string	`json:"name"`
		Rate 	float64	`json:"rate"`
	}

	var resItem []Item

	for _,product := range Products{
		resItem = append(resItem, Item{
			ID : product.ID,
			Name: product.ItemName,
			Rate: product.ProductDetail[0].Rate, 
		})
	}

	c.JSON(http.StatusOK,resItem)
}

func SaveLineSaleItem(c *gin.Context){
	var inventoryReq inventoryRequest
	tokenStr,_ := c.Cookie("JWT-User")
	empId, errTk := helper.DecodeJWT(tokenStr)

	if errTk != nil{
		log.Println("Error while extracting user details")
		c.JSON(http.StatusUnauthorized,gin.H{"error":"Error while extracting user details, please login again"})
		return 
	}


	if err := c.ShouldBindJSON(&inventoryReq); err != nil{
		log.Println("Error while binding Json :",err)
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request"})
		return 
	}

	for _, item := range inventoryReq.Items{

		var product models.Product
		db.Db.Preload("ProductDetail").Where("id = ?",item.ItemID).First(&product)
		// Update inventory

		if product.ProductDetail[0].Stock < int(item.StockOut){
			log.Println("Product not available")
			c.JSON(http.StatusBadRequest,gin.H{"error":"Provided stock not avilable in stocks, please ensure stock update"})
			return
		}

		if item.StockOut > 0 {
			db.Db.Model(&models.ProductDetail{}).Where("product_id = ?",item.ItemID).Update("stock",gorm.Expr("stock - ?",item.StockOut))
		}

		lineSale := models.LineSale{
			EmpID: empId,
			ProductID: product.ID,
			ItemName: product.ItemName,
			Rate: item.Rate,
			StockIn: int(item.StockIn),
			StockOut: int(item.StockOut),
			Damage: int(item.Damage),
			Status: true,
			Created_at: time.Now(),
		}

		if err := db.Db.Create(&lineSale).Error; err != nil{
			log.Println("Error while updating to db :",err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error while saving updates, please try again later"})
			return 
		}
	}

	c.JSON(http.StatusOK,gin.H{"message":"Items saved successfully"})
}

func EditLineSaleItemPage(c *gin.Context){
	itemId := c.Param("id")
	var LineSaleItem models.LineSale
	session := sessions.Default(c)

	if err := db.Db.Where("id = ?",itemId).First(&LineSaleItem).Error; err != nil{
		session.Set("lineError","Could not find item details.")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-items")
		return 
	}

	c.HTML(http.StatusOK,"linesale_item_edit.html",gin.H{
		"ID":LineSaleItem.ID,
		"Name":LineSaleItem.ItemName,
		"Rate":LineSaleItem.Rate,
		"StockIn":LineSaleItem.StockIn,
		"StockOut":LineSaleItem.StockOut,
		"Damage":LineSaleItem.Damage,
		"Status":LineSaleItem.Status,
	})

}

func EditLineSale(c *gin.Context){
	lineSaleId := c.Param("id")
	var LineSaleItem models.LineSale
	session := sessions.Default(c)

	if err := db.Db.Where("id = ?",lineSaleId).First(&LineSaleItem).Error; err != nil{
		session.Set("lineError","Could not fetch details of line sale item")
		session.Save()
		c.Redirect(http.StatusSeeOther,"lineSale-items")
		return 
	}

	StockOutStr := c.PostForm("stock_out")
	stockOut,_ := strconv.Atoi(StockOutStr)

	tempItem := models.LineSaleEdit{
		LineSaleID: LineSaleItem.ID,
		StockOut: stockOut,
	}

	LineSaleItem.Status = false

	if err := db.Db.Create(&tempItem).Error; err != nil{
		session.Set("lineError","Failed to send request for edit")
		session.Save()
		c.Redirect(http.StatusSeeOther,"lineSale-items")
		return 
	}

	if err := db.Db.Save(&LineSaleItem).Error; err != nil{
		session.Set("lineError","Failed to send request for edit")
		session.Save()
		c.Redirect(http.StatusSeeOther,"lineSale-items")
		return
	}

	c.Redirect(http.StatusSeeOther,"lineSale-items")
}