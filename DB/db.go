package db

import (
	"log"
	"os"
	"retialops/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func DbInit(){
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading env file:",err)
	}

	Db,err = gorm.Open(postgres.Open(os.Getenv("dns")),&gorm.Config{})

	if err != nil{
		log.Fatal("Error loading database",err)
	}

	autoerr := Db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Attendance{},
		&models.ProductDetail{},
		&models.CounterSale{},
		&models.LineSale{},
		&models.Credits{},
		&models.Batch{},
		&models.BatchProduction{},
	)

	if autoerr != nil{
		log.Fatal("Migration failed",autoerr)
	}
}