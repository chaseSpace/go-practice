package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAesExample(t *testing.T) {
	//key := []byte(GenRandAesKey())
	key := []byte("fHTLwrrBbICHH/hZ")
	plaintext := []byte(`{"appid":"2021003188669437","application":"shizhen"}`)

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

	/*
		2021003189663400    shizhen20232023@126.com
		CooperateDescKey（时针-支付宝通道1-绑定应用）
		OfAhsH9dOkA/yia+x43H3KTla1QifcBCeZRN0XBQM+ae6GVLsdTo0JZi1hU/Lp0uPfaNA3U5lxjQkt1C70wN81pdPqs=

		2021003188669437	shizhen20232023@163.com
		CooperateDescKey（时针-支付宝通道2-绑定应用）
		0kP/Q7DcV79Fy19ObxOTjHXCi+Hz7I6DMW5X1dT9SqyW26dFSSqWaUQO4PkpjbqdBWnCyZ7PtJe1bE3dIPc6E+SDhl4=
	*/
}
