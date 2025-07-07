package helper

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

func DecodeJWT(tokenStr string) (string,error) {
	secretKey := os.Getenv("secret")
	jwtsecret := []byte(secretKey)

	token, err := jwt.Parse(tokenStr,func(token *jwt.Token)(interface{},error){
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtsecret,nil 
	})

	if err != nil || !token.Valid{
		return "",fmt.Errorf("failed to get value from JWT")
	}

	claims, ok := token.Claims.(jwt.MapClaims);
	if !ok{
		return "",fmt.Errorf("error while getting claims")
	}

	return claims["EmpID"].(string),nil
}