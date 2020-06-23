package util

import (
	"bytes"
	"crypto/des"
	"errors"
	"fmt"
)

func pkcs5Pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPad(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func DesECBEncrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key[:des.BlockSize])
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	bs := block.BlockSize()
	src = pkcs5Pad(src, bs)
	if len(src)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func DesECBDecrypt(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key[:des.BlockSize])
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = pkcs5UnPad(out)
	return out, nil
}
