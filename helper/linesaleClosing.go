package helper

import (
	db "retialops/DB"
	"retialops/models"
	"time"
)

func LineSaleConditions(tokenStr string) (bool, bool) {
	_,superUser, _ := DecodeJWT(tokenStr)

	var todaySts bool
	var count int64

	today := time.Now().Format("2006-01-02")
		
	db.Db.Model(&models.LineSaleClosing{}).Where("created_at BETWEEN ? AND ? ",today+" 00:00:00",today+" 23:59:59").Count(&count)

	if count == 0 {
		todaySts = false 
	}else{
		todaySts = true 
	}

	return superUser,todaySts
}