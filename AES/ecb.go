package AES

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"errors"
)

// 电码本模式（Electronic Codebook Book (ECB)）
type ECB struct {
	Key []byte
}

func New(key string) *ECB {
	err := checkKeyLength(key)
	if err != nil {
		return nil
	}
	return &ECB{Key: []byte(key)}
}

func (e *ECB) Encrypt(data []byte) string {
	// 创建密码组，长度只能是16、24、32字节
	block, _ := aes.NewCipher(e.Key)
	// 获取密钥长度
	blockSize := block.BlockSize()
	// 补码
	data = pkcs7Padding(data, blockSize)
	// 创建保存加密变量
	encryptResult := make([]byte, len(data))
	// CEB是把整个明文分成若干段相同的小段，然后对每一小段进行加密
	for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(encryptResult[bs:be], data[bs:be])
	}
	return base64.StdEncoding.EncodeToString(encryptResult)
}

func (e *ECB) Decrypt(data string) string {
	// 反解密码base64
	originByte, _ := base64.StdEncoding.DecodeString(data)
	// 创建密码组，长度只能是16、24、32字节
	block, _ := aes.NewCipher(e.Key)
	// 获取密钥长度
	blockSize := block.BlockSize()
	// 创建保存解密变量
	decrypted := make([]byte, len(originByte))
	for bs, be := 0, blockSize; bs < len(originByte); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(decrypted[bs:be], originByte[bs:be])
	}
	// 解码
	return string(pkcs7UNPadding(decrypted))
}

func checkKeyLength(key string) error {
	lenMap := map[int]struct{}{
		16: {},
		24: {},
		32: {},
	}
	if _, ok := lenMap[len(key)]; !ok {
		return errors.New("密钥长度不在16,24,32范围内")
	}
	return nil
}

func pkcs7Padding(originByte []byte, blockSize int) []byte {
	// 计算补码长度
	padding := blockSize - len(originByte)%blockSize
	// 生成补码
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	// 追加补码
	return append(originByte, padText...)
}

func pkcs7UNPadding(originDataByte []byte) []byte {
	length := len(originDataByte)
	unPadding := int(originDataByte[length-1])
	return originDataByte[:(length - unPadding)]
}
