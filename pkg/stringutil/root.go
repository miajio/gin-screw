package stringutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const (
	NoBreakSp       = '\u00A0' // no-break space
	NarrowNoBreakSp = '\u202F' // narrow no-break space
	WideNoBreakSp   = '\uFEFF' // wide no-break space
	EnNoBreakSp     = '\u0020' // no-break space
	ZhWideNoBreakSp = '\u3000' // wide no-break space

	LeftToRightEmbedding = '\u202a' // left-to-right embedding
)

// DeSpace delete space in the val
func DeSpace(val string) string {
	result := make([]rune, 0)
	for _, v := range val {
		if v == NoBreakSp || v == NarrowNoBreakSp || v == WideNoBreakSp || v == EnNoBreakSp || v == ZhWideNoBreakSp {
			continue
		}
		result = append(result, v)
	}
	return string(result)
}

// Md5 return md5 string
func Md5(val string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(val)))
}

// Sha256 return sha256 string
func Sha256(val string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(val)))
}

// AesEncrypt return aes encrypt string
func AesEncrypt(v, k string) (string, error) {
	origData, key := []byte(v), []byte(k)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	padding := blockSize - len(origData)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	origData = append(origData, padtext...)

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

// AesDecrypt return aes decrypt string
func AesDecrypt(v, k string) (string, error) {
	crypted, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", err
	}
	key := []byte(k)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	length := len(origData)
	unpadding := int(origData[length-1])
	return string(origData[:(length - unpadding)]), nil
}
