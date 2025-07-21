package main

import (
	"fmt"

	"github.com/capitalone/fpe/ff1"
)

// - FPE 是 Format-Preserving Encryption 的缩写，中文称为“格式保持加密”或“格式保留加密”
// - 它是一种 加密算法，其特点是在加密之后 保持原始数据的格式结构不变

func main() {
	key := []byte("0123456789") // AES密钥
	tweak := []byte("")         // 可选微调数据
	radix := 10                 // 数字 0-9

	cipher, _ := ff1.NewCipher(radix, 7, key, tweak)

	// 加密数字字符串 9527104
	ct, _ := cipher.Encrypt("9527104")
	fmt.Println("加密后：", ct)

	// 解密
	pt, _ := cipher.Decrypt(ct)
	fmt.Println("解密后：", pt)
}
