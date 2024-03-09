package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type Cipher struct {
	gcm cipher.AEAD
}

func NewCipher(shared []byte) (c *Cipher, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(shared)
	if err != nil {
		return nil, err
	}

	var gcm cipher.AEAD
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Cipher{gcm: gcm}, nil
}

func (c Cipher) Encode(src []byte) (dst []byte, err error) {
	nonce := make([]byte, c.gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	dst = c.gcm.Seal(nonce, nonce, src, nil)

	return dst, nil
}

func (c Cipher) Decode(src []byte) (dst []byte, err error) {
	nonceSize := c.gcm.NonceSize()
	nonce, cipherText := src[:nonceSize], src[nonceSize:]

	dst, err = c.gcm.Open(nil, nonce, cipherText, nil)

	return dst, err
}
