package main

import (
	"fmt"
	"time"

	jwt2 "github.com/dgrijalva/jwt-go"

	"toolbox/jwt"
)

func main() {
	secret := "12345"
	claims := jwt.ClaimsData{
		map[string]interface{}{
			"id":   1688,
			"name": "test",
		},
		jwt2.StandardClaims{
			Audience:  "接收JWT者",
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
			Id:        "唯一标识符",
			Issuer:    "发布者",
		},
	}

	j := jwt.New(secret, claims)
	token := j.Create()
	fmt.Println("[create]=", token)

	result, _ := j.Verify(token)
	fmt.Println("[verify]=", result)
}
