package utils

import (
	"fmt"
	"math/rand"
	db "retialops/DB"
	"retialops/models"
	"time"
)

func idGenerator(name string) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return fmt.Sprintf("%s-%d-%d",name,time.Now().UnixNano(),r.Intn(10000))
}

func IDforProduct(name string) string{
	for {
		pID := idGenerator(name)
		var count int64
		db.Db.Model(&models.Product{}).Where("item_id = ?",pID).Count(&count)
		if count == 0{
			return pID
		}
	}
}