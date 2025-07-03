package helper

import (
	"fmt"
	"retialops/models"
	"retialops/utils"
	"time"

	"gorm.io/gorm"
)

func CreateSuperUser(Db *gorm.DB){
	var count int64

	Db.Model(&models.User{}).Where("super_user = ?",true).Count(&count)
	passHash,errp := utils.HashPassword("pass123")
	if errp != nil{
		fmt.Println("Error in hashing admin password :",errp)
	}

	if count == 0 {
		admin := models.User{
			Username: "admin",
			EmpID: "adm1",
			Password: passHash,
			SuperUser: true,
			Created_at: time.Now(),
			Status: "Active",
		}

		if err := Db.Create(&admin).Error; err != nil{
			fmt.Println("Failed to create user :",err)
		}else{
			fmt.Println("Admin created")
		}
	}else{
		fmt.Println("Admin user already exist")
	}
}