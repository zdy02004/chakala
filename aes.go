package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func Encrypt(pass string, aeskey string) (string, error) {

	xpass, err := AesEncrypt([]byte(pass), []byte(aeskey))
	if err != nil {
		return "", err
	}

	pass64 := base64.StdEncoding.EncodeToString(xpass)
	log.Printf("after AesEncrypt:%v\n", pass64)

//	bytesPass, err := base64.StdEncoding.DecodeString(pass64)
//	if err != nil {
//		return "", err
//	}
	return string(pass64), err
}

func Decrypt(pass string, aeskey string) (string, error) {

	tpass, err := base64.StdEncoding.DecodeString( pass )
	if err != nil {
		log.Println(err)
		return "", err
	}

      bytesPass, err1 :=AesDecrypt(tpass,[]byte(aeskey) )
      if err != nil {
              return "", err1
      }
	return string(bytesPass), err
}
