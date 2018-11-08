package encoding

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

var (
	publicKey  []byte
	privateKey []byte
)

// Initialize public and private key.
func InitKey(pubPemPath, priPemPath string) (error) {
	pubKey, err := ioutil.ReadFile(pubPemPath)
	if err != nil {
		os.Exit(-1)
	}
	priKey, err := ioutil.ReadFile(priPemPath)
	if err != nil {
		os.Exit(-1)
	}
	publicKey = pubKey
	privateKey = priKey
	return err
}

// Encrypt plain text by public key, default way is PKCS#1 v1.5.
func RsaEncrypt(plainText []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	cipherData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, plainText)
	return cipherData, err
}

// Decrypt cipher text by private key, default way is PKCS#1 v1.5.
func RsaDecrypt(cipherText []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, pri, cipherText)
	return plainText, err
}

// Sign for md5 data.
func SignPKCS1v15(src, privateKey []byte, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	var err error
	var block *pem.Block
	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
}

func VerifyPKCS1v15(src, sig, publicKey []byte, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	var err error
	var block *pem.Block
	block, _ = pem.Decode(publicKey)
	if block == nil {
		return errors.New("public key error")
	}
	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	var pub = pubInterface.(*rsa.PublicKey)
	return rsa.VerifyPKCS1v15(pub, hash, hashed, sig)
}
