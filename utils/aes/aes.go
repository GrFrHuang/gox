// The AES CBC model using block encryption and decryption of each plaintext block encryption.
// Plain text block length is 128bit, the length of key is 128 or 256bit.
// Padding used to fill the last piece of the 16byte block of plaintext to be encrypted.
// Decrypted using the same padding need to find the last piece of real data in length.
// Most padding modes are PKCS5, PKCS7, NOPADDING.

// For encryption, it should include: encryption secret key length, secret key IV value, encryption mode, padding mode.
// For decryption, it should include: decryption key length, secret key IV value, decryption mode, padding mode.

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const (
	IV = "GreatFriendHuang"
)

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, paddingText...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}

func Encrypt(plainText []byte, key []byte) ([]byte, error) {
	// New a cipher block for Padding encrypted text block.
	// Initialization vector (IV) key's length Support 16bit or 32bit.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	plainText = PKCS5Padding(plainText, blockSize)
	iv := []byte(IV)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

func Decrypt(cipherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, errors.New("cipher text too short")
	}
	iv := []byte(IV)
	if len(cipherText)%blockSize != 0 {
		return nil, errors.New("cipher text is not a multiple of the block size")
	}
	blockModel := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plainText, cipherText)
	plainText = PKCS5UnPadding(plainText)

	return plainText, nil
}
