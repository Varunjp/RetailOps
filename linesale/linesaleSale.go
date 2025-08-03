package linesale

import (
	"log"
	"math"
	"net/http"
	db "retialops/DB"
	"retialops/helper"
	"retialops/models"
	responsemodel "retialops/models/responseModel"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LinesaleSalesPage(c *gin.Context) {

	var Vehicles []models.Vehicle
	session := sessions.Default(c)
	tokenStr,_ := c.Cookie("JWT-User")
	superUser,todayStatus := helper.LineSaleConditions(tokenStr)
	var saleCloses []responsemodel.LineSaleClosing
	slerr := session.Get("saleError")
	message := session.Get("message")

	if err := db.Db.Find(&Vehicles).Error; err != nil{
		c.HTML(http.StatusInternalServerError,"linesale_sales.html",gin.H{"error":"Could not get vehicle details,please try again later"})
		return 
	}

	// pagination for super user table
	pageStr := c.DefaultQuery("page","1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1{
		page = 1
	}
	const pageSize = 10
	var total int64
	db.Db.Model(&models.LineSaleClosing{}).Count(&total)
	totalPages := int(math.Ceil(float64(total)/float64(pageSize)))
	if page > totalPages && totalPages != 0 {
		page = totalPages
	}
	offset := (page - 1)*pageSize

	if superUser {
		var LineSalesClosings []models.LineSaleClosing
		if err := db.Db.Limit(pageSize).Offset(offset).Find(&LineSalesClosings).Error; err != nil{
			log.Println(err)
		}

		for _,closing := range LineSalesClosings{
			var vehicle models.Vehicle
			db.Db.Where("id = ?",closing.Vehicle).First(&vehicle)
			saleCloses = append(saleCloses, responsemodel.LineSaleClosing{
				ID: closing.ID,
				Date: closing.Created_at.Format("2006-01-02"),
				Vehicle: vehicle.Name,
				Sale_amount: closing.SaleAmount,
				Discount: closing.Discount,
				Actual_sale: closing.ActualSale,
				Cash: closing.Cash,
				Account: closing.AccountPayment,
			})
		}
	}

	// pagination conditions
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

	if slerr != nil{
		session.Delete("saleError")
		session.Save()
	}
	if message != nil{
		session.Delete("message")
		session.Save()
	}

	if superUser {
		c.HTML(http.StatusOK,"linesale_sales.html",gin.H{
			"Vehicles":Vehicles,
			"error":slerr,
			"message":message,
			"superUser":superUser,
			"linesaleclosing":saleCloses,
			"todayStatus":todayStatus,
			"CurrentPage":page,
			"TotalPages":totalPages,
			"Pages":pages,
		})
		return 
	}

	c.HTML(http.StatusOK,"linesale_sales.html",gin.H{
		"Vehicles":Vehicles,
		"error":slerr,
		"message":message,
		"superUser":superUser,
		"todayStatus":todayStatus,
		"CurrentPage":page,
		"TotalPages":totalPages,
		"Pages":pages,
	})
}

func LineSaleTotal(c *gin.Context){
	vID := c.Query("vehicle_id")
	today := time.Now().Format("2006-01-02")
	var LineSale []models.LineSale
	session := sessions.Default(c)

	if err := db.Db.Where("vehicle = ? AND created_at BETWEEN ? AND ? ",vID,today+" 00:00:00",today+" 23:59:59").Find(&LineSale).Error; err != nil{
		log.Println("Error while fetching linesale :",err)
		session.Set("saleError","Could not get today's sale details, please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-sales")
		return 
	}

	var totalSales float64

	for _,item := range LineSale {
		var vyapar models.VyaparSale

		if err := db.Db.Where("sale_id = ?",item.ID).First(&vyapar).Error; err != nil{
			log.Println("Error while getting details :",err)
			session.Set("saleError","Error while getting detail, please try again later")
			session.Save()
			c.Redirect(http.StatusSeeOther,"/lineSale-sales")
			return
		}

		totalSales += item.Rate * float64(vyapar.StockOut)
	}

	c.JSON(http.StatusOK,gin.H{"totalAmount":totalSales})
}

func LineSaleSubmit(c *gin.Context){
	
	tokenStr,_ := c.Cookie("JWT-User")
	empId,_,_ := helper.DecodeJWT(tokenStr)
	session := sessions.Default(c)
	vehicleId := c.PostForm("vehicle_id")
	sale_amount,_ := strconv.ParseFloat(c.PostForm("sale_amount"),64)
	discount,_ := strconv.ParseFloat(c.PostForm("discount"),64)
	actual_sale,_ := strconv.ParseFloat(c.PostForm("actual_sale"),64)
	cash,_ := strconv.ParseFloat(c.PostForm("cash"),64)
	account,_ := strconv.ParseFloat(c.PostForm("account"),64)
	var existSalesClose models.LineSaleClosing

	var balance float64
	if actual_sale - (cash + account) <= 0 {
		balance = 0
	}else{
		balance = actual_sale - (cash + account)
	}

	balanceErr := helper.UpdateCreditBalance("lineSale",actual_sale,(cash+account))

	if sale_amount <= 0 || actual_sale <= 0 || cash < 0 || account < 0 {
		session.Set("saleError","Values cannot be zero or less, Please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-sales")
		return 
	}

	if balanceErr != nil{
		session.Set("saleError","Failed to update credit balance, Please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-sales")
		return 
	}

	today := time.Now().Format("2006-01-02")

	if err := db.Db.Where("vehicle = ? AND created_at BETWEEN ? AND ? ",vehicleId,today+" 00:00:00",today+" 23:59:59").First(&existSalesClose).Error; err == nil{
		session.Set("saleError","Sales closing already done for this vehicle today")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-sales")
		return
	}

	lineClosing := models.LineSaleClosing{
		EmpID: empId,
		SaleAmount: sale_amount,
		Discount: discount,
		ActualSale: actual_sale,
		Vehicle: vehicleId,
		Cash: cash,
		AccountPayment: account,
		Balance: balance,
		Created_at: time.Now(),
	}

	if err := db.Db.Create(&lineClosing).Error; err != nil {
		session.Set("saleError","Failed to save line sale closing details, Please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-sales")
		return 
	}

	c.Redirect(http.StatusFound,"/lineSale-sales")
}