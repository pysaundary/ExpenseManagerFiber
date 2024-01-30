package utilities

import (
	"expensesManage/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTToken(IsUser string) (string, error) {
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    IsUser,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	return claim.SignedString([]byte(config.SECRET_KEY))
}

func ParseJWT(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		return "", nil
	}
	claim := token.Claims.(*jwt.StandardClaims)
	return claim.Issuer, nil
}

func GenerateForgetPassWordJWTToken(IsUser string) (string, error) {
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    IsUser,
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})
	return claim.SignedString([]byte(config.SECRET_KEY))
}
