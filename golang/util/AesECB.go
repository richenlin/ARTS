package util

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
)

// 加密
func AesECBEncrypt(plaintext []byte, key string) []byte {
	cipher, err := aes.NewCipher([]byte(key[:aes.BlockSize]))
	if err != nil {
		panic(err.Error())
	}

	if len(plaintext)%aes.BlockSize != 0 {
		panic("Need a multiple of the blocksize 16")
	}

	ciphertext := make([]byte, 0)
	text := make([]byte, 16)
	for len(plaintext) > 0 {
		// 每次运算一个block
		cipher.Encrypt(text, plaintext)
		plaintext = plaintext[aes.BlockSize:]
		ciphertext = append(ciphertext, text...)
	}
	return ciphertext
}

// 解密
func AesECBDecrypt(ciphertext string, key string) []byte {
	cipher, err := aes.NewCipher([]byte(key[:aes.BlockSize]))
	if err != nil {
		panic(err.Error())
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("Need a multiple of the blocksize 16")
	}

	plaintext := make([]byte, 0)
	text := make([]byte, 16)
	byteText, _ := hex.DecodeString(ciphertext)
	for len(byteText) > 0 {
		cipher.Decrypt(text, byteText)
		byteText = byteText[aes.BlockSize:]
		plaintext = append(plaintext, text...)
	}
	return plaintext
}

// Padding补全
func AesECBPKCS5Pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesECBPKCS5UnPad(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesECBPKCS7Pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//
func AesECBPKCS7UnPad(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
