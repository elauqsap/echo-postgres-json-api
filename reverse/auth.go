package reverse

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

type (
	// Protected credentials
	Protected struct {
		Key    string `json:"key"`
		Cipher string `json:"cipher"`
	}
	// Creds unprotected
	Creds struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}
)

// Encrypt creates a protected key cipher pair
func Encrypt(key []byte, plain []byte) (*Protected, error) {
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
	return &Protected{base64.StdEncoding.EncodeToString(key), base64.StdEncoding.EncodeToString(enc)}, nil
}

// Reverse performs reverse encryption of the user:pass
func (p *Protected) Reverse() (*Creds, error) {
	c, err := base64.StdEncoding.DecodeString(p.Key)
	if err != nil {
		return nil, err
	}

	d, err := base64.StdEncoding.DecodeString(p.Cipher)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c)
	if err != nil {
		return nil, err
	}

	if len(d) < aes.BlockSize {
		return nil, errors.New("config a is too short")
	}

	iv := d[:aes.BlockSize]
	d = d[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(d, d)
	if ret := strings.Split(string(d), ":"); len(ret) == 2 {
		return &Creds{User: ret[0], Pass: ret[1]}, nil
	}
	return nil, errors.New("invalid structure to reverse")
}
