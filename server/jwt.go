package server

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/korgottt/go-real-world-api/model"
)

var mySecretKey = []byte(os.Getenv("API_SECRET_KEY"))

//GenerateJwtToken create jwt token
func GenerateJwtToken(username string, id int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":   username,
		"exp":    time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":    time.Now().Unix(),
		"userID": id,
	})
	tokenString, err := token.SignedString(mySecretKey)
	if err != nil {
		return ""
	}
	return tokenString
}

// ParseToken extracts auth data from token
func ParseToken(t string) (d model.SingleUserWrap, e error) {
	token, e := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySecretKey), nil
	})
	if e != nil {
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		d.User.UserName = claims["user"].(string)
	}
	return
}
