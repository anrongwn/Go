package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// RSA public key infor
var (
	devicePublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDtF0kBhMdkCUmBHRARIah1HB44
HCvtOMeqXq1tQ30uvLKMsE3vhT21N1J2Pb2AsX9Q9r/x2abotyfo0f/n2F41YuVN
itn/u7FvEKwc0Uc1gDhQjlIfBTX77uiqaskBGlwAZYfpWtWxakdSl9DeodU4BEhq
OKw1F1KXti5uf1qS5wIDAQAB
-----END PUBLIC KEY-----`)

	notDevicePublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC9Sz8LYpQvOx9x10N9p18C+zdY
yqrxayU1QcNHoFN04uv/IZ6GmnTLQR4HOS093bak0/lkK8wPnXkJHWSuGoqciKkm
WVfW/BugGeH8oUpYkniK8kJZQQOE6wQh4D3h7B3IgIcoFonp2z23gfIIabDvpShF
Hh05JbtIK13I/2RupwIDAQAB
-----END PUBLIC KEY-----`)
)

// init : lib init func
func init() {

}

// RsaEncrypt : rsa encrypt
func RsaEncrypt(plaintext []byte) ([]byte, error) {
	//解密pem公钥
	block, _ := pem.Decode(devicePublicKey)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("device public key error")
	}

	//x509 解析pubkey block
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	//publickey 类型断言
	pub := pubInterface.(*rsa.PublicKey)

	//
	return rsa.EncryptPKCS1v15(rand.Reader, pub, plaintext)
}

// Base64URLCode : Base64URL Code
func Base64URLCode(ciphertext []byte) string {
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// Base64Code : Base64 std Code
func Base64Code(ciphertext []byte) string {
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// GetAesKey : 生成指定位数的 key
func GetAesKey(size int32) ([]byte, error) {
	key := make([]byte, size)
	l, err := rand.Read(key)
	if err != nil || l != len(key) {
		return nil, err
	} else {
		return key, nil
	}
}
