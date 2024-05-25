package util

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

func paddingPKCS7(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func unPaddingPKCS7(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return origData
	}
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesDecryptECB(data, key string) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return []byte{}, err
	}
	decrypted := aesDecryptECB(dataByte, []byte(key))
	return decrypted, nil
}

func aesDecryptECB(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return unPaddingPKCS7(decrypted)
}

func AesEncryptECB(data []byte, key string) string {
	block, _ := aes.NewCipher([]byte(key))
	data = paddingPKCS7(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}

	return base64.StdEncoding.EncodeToString(decrypted)
}

func AesEncryptCBC(plainText string, key string, iv string) string {
	data := aesEncryptCBC([]byte(plainText), []byte(key), []byte(iv))

	return base64.StdEncoding.EncodeToString(data)
}

// AES/CBC/PKCS7Padding
func aesEncryptCBC(plaintext []byte, key []byte, iv []byte) []byte {
	// AES
	block, _ := aes.NewCipher(key)
	// PKCS7 padding
	plaintext = paddingPKCS7(plaintext, aes.BlockSize)
	// CBC encrypt
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(plaintext, plaintext)

	return plaintext
}

func AesDecryptCBC(cipherText string, key string, iv string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return []byte{}, err
	}

	dnData, err := aesDecryptCBC(data, []byte(key), []byte(iv))
	if err != nil {
		return []byte{}, err
	}

	return dnData, nil
}

// aesCBCDecrypt AES/CBC/PKCS7Padding
func aesDecryptCBC(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	// AES
	block, _ := aes.NewCipher(key)
	if len(ciphertext)%aes.BlockSize != 0 {
		return []byte{}, errors.New("ciphertext is not a multiple of the block size")
	}
	// CBC decrypt
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	// PKCS7 unPadding
	result := unPaddingPKCS7(ciphertext)
	return result, nil
}

func RsaSignWithSha256(data, key []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}
	privateKeyByte, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ParsePKCS8PrivateKey err:%v", err))
	}
	privateKey := privateKeyByte.(*rsa.PrivateKey)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error from signing err:%v", err))
	}

	return signature, nil
}

func SignatureCodeWithSha256(value, key string) string {
	dataToHash := []byte(fmt.Sprint(value, key))
	hashToValidate := sha256.Sum256(dataToHash)
	return fmt.Sprintf("%x", hashToValidate)
}

func VerifySignatureCodeWithSha256(value, key, signatureValue string) bool {
	dataToHash := []byte(fmt.Sprint(value, key))
	hashToValidate := sha256.Sum256(dataToHash)
	signature := fmt.Sprintf("%x", hashToValidate)
	return signature == signatureValue
}
