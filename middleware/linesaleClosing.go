package middleware

import (
	"net/http"
	db "retialops/DB"
	"retialops/models"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CheckLineSaleClosing() gin.HandlerFunc{
	return func (c *gin.Context){

		var count int64
		today := time.Now().Format("2006-01-02")
		session := sessions.Default(c)
		
		db.Db.Model(&models.LineSale{}).Where("created_at BETWEEN ? AND ? ",today+" 00:00:00",today+" 23:59:59").Count(&count)

		if count == 0{
			session.Set("lineError","Could not load closing, please enter sale items first")
			session.Save()
			c.Redirect(http.StatusSeeOther,"/lineSale-items")
			c.Abort()
			return 
		}

		c.Next()

	}
}