package utils

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// 生成密码
func GenPassword(password string) string {
	hashStr, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashStr)
}

// 验证密码
func CheckPassword(encryptPassword, enablePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(enablePassword))
	return !(err != nil)
}

// MD5加密
func StrMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 生成UUID
func GenUUID() string {
	return uuid.New().String()
}

// 生成短UUID
func GenShortUUID() string {
	//return shortuuid.NewWithNamespace("order")

	//alphabet := "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy="
	//return shortuuid.NewWithAlphabet(alphabet)

	return shortuuid.New()
}

// 生成安全地全球唯一ID
func GenXID() string {
	return xid.New().String()
}


