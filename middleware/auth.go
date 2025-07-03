package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type Claims struct{
	EmpID 		string
	IsSuperUser	bool 
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading env file:", err)
		}

		key := os.Getenv("secret")
		Secret := []byte(key)
		token, err := c.Cookie("JWT-User")

		if err != nil{
			c.HTML(http.StatusUnauthorized,"login.html",gin.H{"error":"Login required"})
			c.Abort()
			return 
		}

		if token == ""{
			c.HTML(http.StatusUnauthorized,"login.html",gin.H{"error":"Please login again"})
			c.Abort()
			return 
		}

		tokenres, err := jwt.ParseWithClaims(token,&Claims{},func(token *jwt.Token) (interface{}, error) {
			if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
				return nil, fmt.Errorf("unexpected siging method: %v",token.Header["alg"])
			}
			return Secret,nil 
		})

		if err != nil || !tokenres.Valid{
			c.HTML(http.StatusUnauthorized,"login.html",gin.H{"error":"Token expired please login again"})
			c.Abort()
			return 
		}

		c.Next()
	}
}