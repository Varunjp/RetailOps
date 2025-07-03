package helper

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type Claims struct{
	EmpID 		string
	IsSuperUser	bool 
	jwt.StandardClaims
}

func CreateToken(empid string, issuper bool) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file:", err)
	}

	Secret := os.Getenv("secret")
	
	claims := Claims{
		EmpID: empid,
		IsSuperUser: issuper,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 	time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:	time.Now().Unix(),
			Issuer: "Samrudhi Internationals",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	
	return token.SignedString(Secret)
}