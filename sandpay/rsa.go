package sandpay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

const signField = "sign"

var sandSignKey = map[string]int{
	"version":       1,
	"mer_no":        1,
	"mer_order_no":  1,
	"create_time":   1,
	"order_amt":     1,
	"notify_url":    1,
	"return_url":    1,
	"create_ip":     1,
	"pay_extra":     1,
	"accsplit_flag": 1,
	"sign_type":     1,
	"store_id":      1,
}

func signWithPKCS1v15(param url.Values, privateKey *rsa.PrivateKey, hash crypto.Hash) (s string, err error) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for key := range param {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 && sandSignKey[key] == 1 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	var src = strings.Join(pList, "&")

	sig, err := RSASignWithKey([]byte(src), privateKey, hash)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

func verifySign(data url.Values, key *rsa.PublicKey) (ok bool, err error) {
	sign := data.Get(signField)

	var keys = make([]string, 0, 0)
	for k := range data {
		if k == signField {
			continue
		}
		keys = append(keys, k)
	}

	sort.Strings(keys)
	var buf strings.Builder

	for _, k := range keys {
		vs := data[k]
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	s := buf.String()
	return verify([]byte(s), sign, key)
}

func verify(data []byte, sign string, key *rsa.PublicKey) (ok bool, err error) {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	if err = RSAVerifyWithKey(data, signBytes, key, crypto.SHA1); err != nil {
		return false, err
	}
	return true, nil
}

func RSASignWithKey(data []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(data)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func RSAVerifyWithKey(data, sig []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(data)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sig)
}
