package linesale

import (
	"log"
	"net/http"
	db "retialops/DB"
	"retialops/models"
	"retialops/models/getitems"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LineSaleClosingPage(c *gin.Context) {
	var Vehicles []models.Vehicle
	session := sessions.Default(c)

	slcerr := session.Get("clserr")
	msg := session.Get("message")

	if err := db.Db.Find(&Vehicles).Error; err != nil{
		log.Println("Failed to load vehicle details :",err)
		c.HTML(http.StatusInternalServerError,"linesale_closing.html",gin.H{"error":"Failed to get vehicles, please try again later"})
		return 
	}

	if slcerr != nil{
		session.Delete("clserr")
		session.Save()
	}
	if msg != nil{
		session.Delete("message")
		session.Save()
	}

	c.HTML(http.StatusOK,"linesale_closing.html",gin.H{
		"Vehicle":Vehicles,
		"error":slcerr,
		"message":msg,
	})
}

func LineSaleClosingDetails(c *gin.Context){
	vId := c.Query("vehicle_id")
	today := time.Now().Format("2006-01-02")
	var lineSales []getitems.Linetable
	var sales 	getitems.Revenue
	var expenses []getitems.Expenses
	var TodaySalesItem []models.LineSale
	var TodaySale models.LineSaleClosing
	var Todayexpense 	[]models.LineSaleExpenses
	session := sessions.Default(c)

	if err := db.Db.Where("vehicle = ? AND created_at BETWEEN ? AND ? ",vId,today+" 00:00:00",today+" 23:59:59").Find(&TodaySalesItem).Error; err != nil{
		log.Println("Failed to fetch items :",err)
		session.Set("clserr","Failed to fetch item details")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-closing")
		return 
	}

	if err := db.Db.Where("vehicle = ? AND created_at BETWEEN ? AND ? ",vId,today+" 00:00:00",today+" 23:59:59").First(&TodaySale).Error; err != nil{
		log.Println("Failed to fetch sale details :",err)
		session.Set("clserr","Failed to fetch sale details")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-closing")
		return
	}

	if err := db.Db.Where("line_sale_id = ?",TodaySale.ID).Find(&Todayexpense).Error; err != nil{
		log.Println("Failed to fetch expense details :",err)
		session.Set("clserr","Failed to fetch expense details")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/lineSale-closing")
		return
	}

	for _,item := range TodaySalesItem{
		lineSales = append(lineSales, getitems.Linetable{
			Name: item.ItemName,
			StockIn: item.StockIn,
			StockOut: item.StockOut,
			Rate: item.Rate,
		})
	}

	sales.ActualSale = TodaySale.ActualSale
	sales.Discount = TodaySale.Discount
	sales.TotalAmount = TodaySale.SaleAmount
	sales.Cash = TodaySale.Cash
	sales.Account = TodaySale.AccountPayment

	for _, exp := range Todayexpense{
		expenses = append(expenses, getitems.Expenses{
			Item: exp.Type,
			Amount: exp.Amount,
		})
	}

	c.JSON(http.StatusOK,gin.H{
		"line_sales": lineSales,
        "revenue": sales,
        "expenses": expenses,
	})
}