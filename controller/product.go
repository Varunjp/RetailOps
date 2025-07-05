package controller

import (
	"math"
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
	var InactiveProducts []models.Product
	
	pageStr := c.DefaultQuery("page","1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1{
		page = 1
	} 

	const pageSize = 10
	var total int64
	db.Db.Model(&models.Product{}).Where("status = ?","Active").Count(&total)

	totalPages := int(math.Ceil(float64(total)/float64(pageSize)))

	if page > totalPages && totalPages != 0 {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	if err := db.Db.Preload("ProductDetail").Where("status = ?","Active").Order("id DESC").Limit(pageSize).Offset(offset).Find(&Products).Error; err != nil{
		c.HTML(http.StatusNotFound,"products.html",gin.H{"error":"Failed to return products or No products added"})
		return 
	}

	if err := db.Db.Preload("ProductDetail").Where("status != ?","Active").Order("id DESC").Limit(pageSize).Offset(offset).Find(&InactiveProducts).Error; err != nil{
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
	var inactiveResponseProduct []response

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

	for _, product := range InactiveProducts{
		rate := product.ProductDetail[0].Rate
		stock := product.ProductDetail[0].Stock
		inactiveResponseProduct = append(inactiveResponseProduct, response{
			ID: product.ID,
			Name: product.ItemName,
			Rate: rate,
			Stock: stock,
		})
	}
	

	session := sessions.Default(c)
	errmsg := session.Get("product")
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

	if errmsg != nil{
		session.Delete("product")
		session.Save()
		c.HTML(http.StatusOK,"products.html",gin.H{
			"products":responseProduct,
			"inactproducts":inactiveResponseProduct,
			"error":errmsg,
			"CurrentPage": page,
			"TotalPages":totalPages,
			"Pages": pages,
		})
		return
	}

	c.HTML(http.StatusOK,"products.html",gin.H{
		"products":responseProduct,
		"inactproducts":inactiveResponseProduct,
		"error":errmsg,
		"CurrentPage": page,
		"TotalPages":totalPages,
		"Pages": pages,
	})
}

func AddProduct(c *gin.Context){
	ProductName := c.PostForm("productName")
	RateStr := c.PostForm("rate")
	StockStr := c.PostForm("stock")
	session := sessions.Default(c)
	var productCheck models.Product

	if strings.TrimSpace(ProductName) == "" || strings.TrimSpace(RateStr) == "" || strings.TrimSpace(StockStr) == ""{
		session.Set("product","Invalid input passed")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	if err := db.Db.Where("item_name ILIKE ?",ProductName).First(&productCheck).Error; err == nil{
		session.Set("product","Product already exist")
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

func EditProduct(c *gin.Context){
	session := sessions.Default(c)
	productId := c.Param("id")
	ProductName := c.PostForm("productName")
	RateStr := c.PostForm("rate")
	StockStr := c.PostForm("stock")
	status := c.PostForm("status")
	var product models.Product

	if strings.TrimSpace(ProductName) == "" || strings.TrimSpace(RateStr) == "" || strings.TrimSpace(StockStr) == ""{
		session.Set("product","Invalid input passed")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	rate,_ := strconv.ParseFloat(RateStr,64)
	stock,_ := strconv.Atoi(StockStr)
	var productDetails models.ProductDetail

	if err := db.Db.Preload("ProductDetail").Where("id = ?",productId).First(&product).Error; err != nil{
		session.Set("product","Could not find product details, please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	if err := db.Db.Where("product_id = ?",product.ID).First(&productDetails).Error; err != nil{
		session.Set("product","Could not find product details, please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	if rate < 1 {
		session.Set("product","Product cannot be updated to rate less than zero")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	if rate != 0 {
		productDetails.Rate = rate
	}

	if stock != 0 {
		productDetails.Stock+=stock
	}
	
	if ProductName != product.ItemName{
		product.ItemName = ProductName
	}

	if status != product.Status{
		product.Status = status
	}

	if err := db.Db.Save(&product).Error; err != nil{
		session.Set("product","Failed to save updates, please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}

	if err := db.Db.Save(&productDetails).Error; err != nil{
		session.Set("product","Failed to save updates, please try again later")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/user/products")
		return 
	}
	c.Redirect(http.StatusSeeOther,"/user/products")
}