package helper

import (
	db "retialops/DB"
	"retialops/models"
	"time"

	"gorm.io/gorm"
)

func UpdateCreditBalance(typesale string, saleamount, paidamount float64) error {

	var Credit models.Credits

	if err := db.Db.Where("type = ?",typesale).First(&Credit).Error; err != nil{
		if err == gorm.ErrRecordNotFound{
			creditAccount := models.Credits{
				Type: typesale,
				Created_at: time.Now(),
			}

			if err := db.Db.Create(&creditAccount).Error; err != nil{
				return err 
			}

			balance := saleamount - paidamount
			creditAccount.Balance = creditAccount.Balance + balance

			if err := db.Db.Save(&creditAccount).Error; err != nil{
				return err 
			}
			return nil 
			
		}else{
			return err
		}
		 
	}

	balance := saleamount - paidamount
	Credit.Balance = Credit.Balance + balance

	if err := db.Db.Save(&Credit).Error; err != nil{
		return err 
	}

	return nil 
}