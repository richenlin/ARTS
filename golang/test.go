/**
 * @Author: richen
 * @Date: 2017-12-08 16:25:47
 * @Copyright (c) - <richenlin(at)gmail.com>
 * @Last Modified by: richen
 * @Last Modified time: 2017-12-08 19:38:06
 */
package main

import (
	"crypto/aes"
	"demo/util"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func test() {
	// testAesECB()
	testAesCBC()
	// testDesCBC()
	// testDesECB()
	// testNeo()
}

func testAesECB() {
	fmt.Println("testAesECB ---------------------------------")
	txt := "201711270101000000001"
	key := "4uoRHo0TC62DyLPh7QYlWA=="

	heKey, _ := base64.StdEncoding.DecodeString(key)
	// byteData := aesECB.Encrypt(aesECB.PKCS7Pad([]byte(txt), 128), string(heKey))
	byteData := util.AesECBEncrypt(util.AesECBPKCS7Pad([]byte(txt), aes.BlockSize), string(heKey))
	cipherText := hex.EncodeToString(byteData)
	fmt.Println(cipherText)
	fmt.Printf("base64 ciphertext: %s \n", base64.StdEncoding.EncodeToString(byteData))
	plainText := util.AesECBPKCS7UnPad(util.AesECBDecrypt(cipherText, string(heKey)))
	fmt.Printf("plaintext: %s \n", plainText)

}

func testAesCBC() {
	fmt.Println("testAesCBC ---------------------------------")
	txt := testJSON()
	key := "12345678901234567890123456789012"
	fmt.Println(txt)
	// heKey, _ := base64.StdEncoding.DecodeString(key)
	byteData, _ := util.AesCBCEncrypt([]byte(txt), []byte(key))
	cipherText := hex.EncodeToString(byteData)
	fmt.Println(cipherText)
	fmt.Printf("base64 ciphertext: %s \n", base64.StdEncoding.EncodeToString(byteData))
	plainText, _ := util.AesCBCDecrypt(byteData, []byte(key))
	fmt.Printf("plaintext: %s \n", plainText)
}

func testDesCBC() {
	fmt.Println("testDesCBC ---------------------------------")
	txt := "201711270101000000001"
	key := "12345678"

	result, _ := util.DesCBCEncrypt([]byte(txt), []byte(key))
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	origData, _ := util.DesCBCDecrypt(result, []byte(key))
	fmt.Println(string(origData))
}

func testJSON() string {
	fmt.Println("testJSON ---------------------------------")
	type Info struct {
		Name   string `json:"name"`
		Mobile string `json:"mobile"`
		Cardno string `json:"idcard"`
	}

	info := Info{
		Name:   "张三",
		Mobile: "13800138000",
		Cardno: "101111111111111111",
	}

	txt, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	res := string(txt)
	fmt.Println(res)
	return res
}

func testDesECB() {
	fmt.Println("testDesECB ---------------------------------")
	txt := "201711270101000000001"
	key := "4uoRHo0TC62DyLPh7QYlWA=="

	heKey, _ := base64.StdEncoding.DecodeString(key)
	result, _ := util.DesECBEncrypt([]byte(txt), heKey)
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	origData, _ := util.DesECBDecrypt(result, heKey)
	fmt.Println(string(origData))
}

func testNeo() {
	// app := neo.App()

	// app.Get("/", func(ctx *neo.Ctx) (int, error) {
	// 	token := "7d2758efad43f37ff1957018429c3946"
	// 	type Info struct {
	// 		token       string
	// 		Name        string
	// 		Cardno      string `json:"cardno"`
	// 		Mobile      string `json:"mobile"`
	// 		CompanyName string `json:"companyName"`
	// 		Salary      string `json:"salary"`
	// 	}
	// 	info := Info{
	// 		token:       token,
	// 		Name:        "木木",
	// 		Cardno:      "110000000000000000",
	// 		Mobile:      "",
	// 		CompanyName: "北京xxxx有限公司",
	// 		Salary:      "1000",
	// 	}
	// 	return 200, ctx.Res.Json(info)
	// })

	// app.Use(func(ctx *neo.Ctx, next neo.Next) {
	// 	start := time.Now().UnixNano()
	// 	fmt.Printf("--> [Req] %s to %s", ctx.Req.Method, ctx.Req.URL.Path)

	// 	next()

	// 	elapsed := int64(time.Now().UnixNano()-start) / 1000
	// 	fmt.Printf("<-- [Res] (%d) %s to %s Took %vµs \n", ctx.Res.Status, ctx.Req.Method, ctx.Req.URL.Path, elapsed)
	// })

	// app.Start()
}
