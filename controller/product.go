package controller

import (
	"net/http"
	db "retialops/DB"
	"retialops/models"
	"retialops/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ProductPage(c *gin.Context) {
	var Products []models.Product

	if err := db.Db.Preload("ProductDetail").Where("status = ?","Active").Find(&Products).Error; err != nil{
		c.HTML(http.StatusNotFound,"products.html",gin.H{"error":"Failed to return products or No products added"})
		return 
	}

	type response struct{
		ID			uint
		Name		string 
		Rate		float64
		Stock 		int 
	}

	var responseProduct []response

	for _, product := range Products{
		rate := product.ProductDetail[0].Rate
		stock := product.ProductDetail[0].Stock
		responseProduct = append(responseProduct, response{
			ID: product.ID,
			Name: product.ItemName,
			Rate: rate,
			Stock: stock,
		})
	}

	session := sessions.Default(c)
	errmsg := session.Get("product")

	if errmsg != nil{
		session.Delete("product")
		session.Save()
		c.HTML(http.StatusOK,"products.html",gin.H{
			"products":responseProduct,
			"error":errmsg,
		})
		return
	}

	c.HTML(http.StatusOK,"products.html",gin.H{
		"products":responseProduct,
	})
}

func AddProduct(c *gin.Context){
	ProductName := c.PostForm("productName")
	RateStr := c.PostForm("rate")
	StockStr := c.PostForm("stock")
	session := sessions.Default(c)

	if strings.TrimSpace(ProductName) == "" || strings.TrimSpace(RateStr) == "" || strings.TrimSpace(StockStr) == ""{
		session.Set("product","Invalid input passed")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	rate,_ := strconv.ParseFloat(RateStr,64)
	stock,_ := strconv.Atoi(StockStr)

	if rate < 1 || stock < 1{
		session.Set("product","Rate or stock cannot be below zero")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	productID := utils.IDforProduct(ProductName)
	timeNow := time.Now()
	date := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, timeNow.Location())

	product := models.Product{
		ItemName: ProductName,
		ItemID: productID,
		Status: "Active",
		Created_at: date,
	}

	if err := db.Db.Create(&product).Error; err != nil{
		session.Set("product","Failed to create product, please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	productDetails := models.ProductDetail{
		ProductID: product.ID,
		Rate: rate,
		Stock: stock,
	}

	if err := db.Db.Create(&productDetails).Error; err != nil{
		session.Set("product","Failed to create product, please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	c.Redirect(http.StatusSeeOther,"/user/products")
}

func EditProductPage(c *gin.Context){
	productID := c.Param("id")
	var product models.Product
	session := sessions.Default(c)

	if err := db.Db.Preload("ProductDetail").Where("id = ?",productID).First(&product).Error; err != nil{
		session.Set("product","Failed to retrive product details")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	c.HTML(http.StatusOK,"product_edit.html",gin.H{
		"ID":productID,
		"Name":product.ItemName,
		"Rate": product.ProductDetail[0].Rate,
		"Stock": product.ProductDetail[0].Stock,
		"Status": product.Status,
	})
}