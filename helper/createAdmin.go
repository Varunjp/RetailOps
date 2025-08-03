package helper

import (
	"fmt"
	"log"
	"os"
	"retialops/models"
	"retialops/utils"
	"time"

	"gorm.io/gorm"
)

func CreateSuperUser(Db *gorm.DB){
	var count int64

	Db.Model(&models.User{}).Where("super_user = ?",true).Count(&count)
	password := os.Getenv("ADMIN_PASSWORD")
	passHash,errp := utils.HashPassword(password)
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

func CreateNormalUser(Db *gorm.DB){
	var User models.User

	if err := Db.Where("emp_id = ?","test").First(&User).Error; err != nil{
		if err == gorm.ErrRecordNotFound{
			hashpass,passerr := utils.HashPassword("pass")
			if passerr != nil{
				log.Println(passerr)
			}
			newUser := models.User{
				Username: "test",
				EmpID: "test",
				Phone: "7438293728",
				Password: hashpass,
				Created_at: time.Now(),
				SuperUser: false,
				Status: "Active",
			}

			if err := Db.Create(&newUser).Error; err != nil{
				log.Println(err)
			}else{
				log.Println("User created successfully")
			}
		}
	}

	log.Println("User already exist")
}