package vehicles

import (
	"log"
	"math"
	"net/http"
	db "retialops/DB"
	"retialops/models"
	responsemodel "retialops/models/responseModel"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VehicleList(c *gin.Context){
	var Vehicles []models.Vehicle
	var resVehicles []responsemodel.Vehicle

	if err := db.Db.Find(&Vehicles).Error; err != nil{
		log.Println("Error while getting vehicle details :",err)
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Could not get vehicle details, please try again later"})
		return 
	}

	for _,veh := range Vehicles{
		resVehicles = append(resVehicles, responsemodel.Vehicle{
			ID: veh.ID,
			Name: veh.Name,
			License: veh.License,
		})
	}

	c.JSON(http.StatusOK,resVehicles)
}

func VehiclesPage(c *gin.Context){
	var Vehicles []models.Vehicle
	session := sessions.Default(c)

	pageStr := c.DefaultQuery("page","1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1{
		page = 1
	} 

	const pageSize = 10
	var total int64
	db.Db.Model(&models.Vehicle{}).Count(&total)
	totalPages := int(math.Ceil(float64(total)/float64(pageSize)))

	if page > totalPages && totalPages != 0 {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	if err := db.Db.Limit(pageSize).Offset(offset).Find(&Vehicles).Error; err != nil{
		if err != gorm.ErrRecordNotFound{
			c.HTML(http.StatusInternalServerError,"vehicles.html",gin.H{"error":"Could not load any vehicle details"})
			return 
		}
	}

	veErr := session.Get("vehicleError")
	message := session.Get("message")

	if veErr != nil{
		session.Delete("vehicleError")
		session.Save()
	}
	if message != nil{
		session.Delete("message")
		session.Save()
	}

	// pagination 
	var pages []int 
	startPage := max(page - 2, 1)
	endPage := min(startPage + 4, totalPages)

	for i := startPage; i <= endPage; i++{
		pages = append(pages, i)
	}

	c.HTML(http.StatusOK,"vehicles.html",gin.H{
		"vehicles":Vehicles,
		"CurrentPage":page,
		"TotalPages":totalPages,
		"Pages":pages,
		"error":veErr,
		"message":message,
	})
}

func AddVehicle(c *gin.Context){
	name := c.PostForm("vehicleName")
	license := c.PostForm("license")
	var CheckVehicle models.Vehicle
	session := sessions.Default(c)

	if err := db.Db.Where("name ILIKE ? AND license ILIKE ?",name,license).First(&CheckVehicle).Error; err == nil{
		log.Println("Vehicle already exist")
		session.Set("vehicleError","Vehicle already exist")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/smrd/vehicles")
		return 
	}

	if strings.TrimSpace(name) == "" || strings.TrimSpace(license) == ""{
		log.Println("Incorrect value passed")
		session.Set("vehicleError","Passed value doesn't satisfy name or license")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/smrd/vehicles")
		return 
	}

	newVehicle := models.Vehicle{
		Name: name,
		License: license,
		Created_at: time.Now(),
	}

	if err := db.Db.Create(&newVehicle).Error; err != nil{
		log.Println("Failed to save new vehicle :",err)
		session.Set("vehicleError","Failed to create new vehicle")
		session.Save()
		c.Redirect(http.StatusSeeOther,"/smrd/vehicles")
		return 
	}

	session.Set("message","Successfully created new vehicle")
	session.Save()
	c.Redirect(http.StatusFound,"/smrd/vehicles")
}