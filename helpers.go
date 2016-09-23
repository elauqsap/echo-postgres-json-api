package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func encrypt(key []byte, plain []byte) (*string, error) {
	ret := new(string)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	enc := make([]byte, aes.BlockSize+len(plain))
	iv := enc[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(enc[aes.BlockSize:], plain)
	tmp := base64.StdEncoding.EncodeToString(enc)
	ret = &tmp
	return ret, nil
}
