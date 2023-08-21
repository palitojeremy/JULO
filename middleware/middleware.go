package middleware

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func SplitToken(headerToken string) string {
	parsedToken := strings.SplitAfter(headerToken, " ")
	tokenString := parsedToken[1]
	return tokenString
}

func AuthenticateToken(tokenString string) error {
	//token check
	_, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil

}

func DecodeToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	if err != nil {
		return token, err
	}
	return token, nil
}