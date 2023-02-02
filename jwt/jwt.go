package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	Secret string
	Claims ClaimsData
}

type ClaimsData struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

func New(secret string, claims ClaimsData) *Jwt {
	return &Jwt{
		Secret: secret,
		Claims: claims,
	}
}

// 创建令牌
func (j *Jwt) Create() string {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, j.Claims)
	result, err := tokenObj.SignedString([]byte(j.Secret))
	if err != nil {
		return ""
	}
	return result
}

// 校验令牌
func (j *Jwt) Verify(token string) (*ClaimsData, interface{}) {
	var c ClaimsData
	tokenObj, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token不正确")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, fmt.Errorf("token已过期")
			} else {
				return nil, fmt.Errorf("token格式错误")
			}
		}
	}

	if tokenObj != nil {
		if key, ok := tokenObj.Claims.(*ClaimsData); ok && tokenObj.Valid {
			return key, nil
		} else {
			return nil, fmt.Errorf("token不正确")
		}
	}
	return nil, fmt.Errorf("token不正确")
}
