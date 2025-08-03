package linesale

import (
	"log"
	"math"
	"net/http"
	db "retialops/DB"
	"retialops/models"
	responsemodel "retialops/models/responseModel"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VyparPage(c *gin.Context){
	session := sessions.Default(c)
	message := session.Get("message")
	vpyerr := session.Get("vyapar_error")

	if message != nil{
		session.Delete("message")
		session.Save()
	}
	if vpyerr != nil{
		session.Delete("vyapar_error")
		session.Save()
	}

	c.HTML(http.StatusOK,"vyapar.html",gin.H{
		"error":vpyerr,
		"message":message,
	})
}

// provide items for selected vehicles
func VyaparStockOut(c *gin.Context){
	var Itemslist []models.LineSale
	var Items []responsemodel.LineSaleStockOut
	vehicle := c.Query("vehicle_id")
	today := time.Now().Format("2006-01-02")
	query := db.Db.Model(&models.LineSale{}).Where("vehicle = ? AND created_at BETWEEN ? AND ? AND status = ?",vehicle,today+" 00:00:00",today+" 23:59:59",true)

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

func VyaparUpdate(c *gin.Context){
	var vypReq responsemodel.VyaparUpdateRequest
	var vypResp []responsemodel.VyaparResponse
	var existVyp models.VyaparSale
	var remains	bool 
	session := sessions.Default(c)

	if err := c.ShouldBindJSON(&vypReq); err != nil{
		log.Println("Error while binding data in vyapar update :",err)
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request received"})
		return 
	}

	if len(vypReq.StockUpdates) == 0{
		log.Println("not items passes")
		c.JSON(http.StatusBadRequest,gin.H{"error":"No items provided"})
		return 
	}

	if err := db.Db.Where("sale_id = ?",vypReq.StockUpdates[0].ItemID).First(&existVyp).Error; err == nil{
		log.Println("Vyapar sale already exists for this item")
		c.JSON(http.StatusConflict,gin.H{"error":"Vyapar sale already exists for this item, cannot update"})
		return 
	}

	for _,item := range vypReq.StockUpdates{
		var linesale models.LineSale
		if err := db.Db.Where("id = ?",item.ItemID).First(&linesale).Error; err != nil{
			log.Println("Failed to fetch item details in vyapar update :",err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not fetch details, please try again later"})
			return 
		}

		if item.StockOut < 0 {
			c.JSON(http.StatusBadRequest,gin.H{"error":"Stock out cannot be less than zero"})
			return 
		}

		balance := linesale.StockIn - linesale.StockOut
		vpyr := item.StockOut

		if balance != vpyr{
			remains = true 
			diff := balance - vpyr
			if diff > 0 {
				vypResp = append(vypResp,responsemodel.VyaparResponse{
					ItemName: linesale.ItemName,
					Difference: diff,
					Status: "Less",
				})
			}else{
				vypResp = append(vypResp,responsemodel.VyaparResponse{
					ItemName: linesale.ItemName,
					Difference: int(math.Abs(float64(diff))),
					Status: "More",
				})
			}
		}
	}

	if remains {
		c.JSON(http.StatusBadRequest,gin.H{
			"mismatches" : vypResp,
		})
		return 
	}

	for _,item := range vypReq.StockUpdates{
		var linesale models.LineSale

		if err := db.Db.Where("id = ?",item.ItemID).First(&linesale).Error; err != nil{
			log.Println("Failed to fetch item details in vyapar update :",err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not fetch item details, please try again later"})
			return 
		}

		vypSale := models.VyaparSale{
			SaleID: linesale.ID,
			StockOut: item.StockOut,
			Created_at: time.Now(),
		}

		if err := db.Db.Create(&vypSale).Error; err != nil{
			log.Println("Could not save sale item :",err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not save item, please try again later"})
			return 
		}

		if err := db.Db.Model(&models.ProductDetail{}).Where("product_id = ?",linesale.ProductID).Update("stock",gorm.Expr("stock - ?",item.StockOut)).Error; err != nil{
			log.Println("Failed to update product details :",err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to update stock update of product"})
			return 
		}
	}
	session.Set("message","Vyapar update successfull")
	session.Save()
	c.Redirect(http.StatusFound,"/lineSale-vypar")	
}