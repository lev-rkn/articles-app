package utils

import "github.com/golang-jwt/jwt/v5"

type extractedData struct {
	UserId int
}

func ParseUserClaims(ctxValue any) *extractedData {
	data := extractedData{}

	claims := ctxValue.(*jwt.Token).Claims.(jwt.MapClaims)
	data.UserId = int(claims["uid"].(float64))

	return &data
}
