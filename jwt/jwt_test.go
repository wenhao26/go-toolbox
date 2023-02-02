package jwt

import (
	"testing"
	"time"

	jwt2 "github.com/dgrijalva/jwt-go"
)

var secret = "12345"
var claims = ClaimsData{
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

func TestCreate(t *testing.T) {
	t.Log(New(secret, claims).Create())
}

func TestVerify(t *testing.T) {
	t.Log(New(secret, claims).Verify("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImlkIjoxNjg4LCJuYW1lIjoidGVzdCJ9LCJhdWQiOiLmjqXmlLZKV1TogIUiLCJleHAiOjE2NzM2MDE4MTAsImp0aSI6IuWUr-S4gOagh-ivhuespiIsImlzcyI6IuWPkeW4g-iAhSJ9.-Swkni_dKrBn4xTY8HKB7x0f2S33E3yvepiy8efNNtA"))
}

func BenchmarkCreateAndVerify(b *testing.B) {
	b.ResetTimer()
	j := New(secret, claims)
	for i := 0; i < b.N; i++ {
		token := j.Create()
		j.Verify(token)
	}
}
