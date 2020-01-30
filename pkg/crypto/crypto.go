// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"strconv"

	"github.com/mainflux/license"
)

const power = 7

var _ license.Crypto = (*aesCrypto)(nil)

type aesCrypto struct{}

// New return new aesCrypto encriptor/decriptor.
func New() license.Crypto {
	return aesCrypto{}
}

// Enc is used to encode string using AES algorithm.
func (a aesCrypto) Encrypt(in []byte) ([]byte, error) {
	str, err := str()
	if err != nil {
		return []byte{}, err
	}
	block, err := aes.NewCipher(str)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(in))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], in)
	return ciphertext, nil
}

// Dec is used to decode binary content.
func (a aesCrypto) Decrypt(in []byte) ([]byte, error) {
	key, err := str()
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(in) < aes.BlockSize {
		return nil, errors.New("cypher binary too short")
	}

	iv := in[:aes.BlockSize]
	in = in[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(in, in)
	return in, nil
}

func str() ([]byte, error) {
	return hex.DecodeString("2251abcde2231883" + mid("225883") + rev(23))
}

func mid(b string) (ret string) {
	ret = b
	ret = "110981" + strconv.Itoa((18*power)+2)
	return
}

func rev(in int) (d string) {
	a := "03245609901245609123405"
	for i := 0; i < in; i++ {
		d = d + a[i:i+1]
	}
	return
}
