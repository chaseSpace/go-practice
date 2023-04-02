package adapayCore

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

func RsaLongEncrypt(content string, msc *MerchSysConfig) (result string) {

	return ""
}

func RsaLongDecrypt(content string, msc *MerchSysConfig) (result string) {

	return ""
}

func RsaSign(content string, msc *MerchSysConfig) (result string, err error) {
	privateKey, err := getPrivateKey(msc)
	if err != nil {
		return "", err
	} else if privateKey == nil {
		return "", errors.New("Error in get private key")
	}

	h := sha1.New()
	h.Write([]byte([]byte(content)))
	hash := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hash[:])
	if err != nil {
		Println("Sign Error: ", err)
		return "", err
	}

	out := base64.StdEncoding.EncodeToString(signature)

	return out, nil
}

func RsaSignVerify(signData string, originalData string, msc *MerchSysConfig) (err error) {
	publicKey, pubk_err := getPublicKey(msc)
	if pubk_err != nil {
		return pubk_err
	}

	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}

	hash := sha1.New()
	hash.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

func getPrivateKey(msc *MerchSysConfig) (interface{}, error) {
	priKeyString := msc.RspPriKey
	if priKeyString == "" {
		return nil, errors.New("Unknow private key data")
	}

	if !strings.Contains(priKeyString, "-----BEGIN PRIVATE KEY-----") {
		priKeyString = "-----BEGIN PRIVATE KEY-----\n" + priKeyString
	}
	if !strings.Contains(priKeyString, "-----END PRIVATE KEY-----") {
		priKeyString = priKeyString + "\n-----END PRIVATE KEY-----"
	}

	defer func() {
		err := recover()
		if err != nil {
			Println("Error in get private key: ", err)
		}
	}()
	keyByts, _ := pem.Decode([]byte(priKeyString))

	privateKey, err := x509.ParsePKCS8PrivateKey(keyByts.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func getPublicKey(msc *MerchSysConfig) (interface{}, error) {
	publicString := msc.RspPubKey
	if publicString == "" {
		return "", errors.New("Unknow public key data")
	}

	if !strings.Contains(publicString, "-----BEGIN PUBLIC KEY-----") {
		publicString = "-----BEGIN PUBLIC KEY-----\n" + publicString
	}
	if !strings.Contains(publicString, "-----END PUBLIC KEY-----") {
		publicString = publicString + "\n-----END PUBLIC KEY-----"
	}

	defer func() {
		err := recover()
		if err != nil {
			Println("Error in get public key: ", err)
		}
	}()
	keyByts, _ := pem.Decode([]byte(publicString))

	publicKey, err := x509.ParsePKIXPublicKey(keyByts.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
