package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAesExample(t *testing.T) {
	//key := []byte(GenRandAesKey())
	key := []byte("fHTLwrrBbICHH/hZ")
	plaintext := []byte(`{"appid":"app_e378030b-fdd1-4756-b4d9-664abc0c11a3"}`)

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		panic(err)
	}
	b64Cipher := base64.StdEncoding.EncodeToString(ciphertext)
	fmt.Printf("ciphertext.base64=%s\n", b64Cipher)
	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("text=%s\n", decrypted)

	//CooperateDescKey（绑定sand）：C/lo5O6OL5emYxLME/Ky4ehN2SJI84csTv5gsXO0tPT4EGeTZ+whWxo=
	//CooperateDescKey（绑定汇付）：b/Ubr5TwiLJAOhMaF+L0XnR8gA9y25OrT8wlXrzO7gRsDXq0XyhWoqbIokOYno1A0lQHsS6lSZrcqyAwZ+1KjvqhWX4=
}
