package tencent_svc

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

func IMAuditText(text string, appId int, key, identifier string) (*AuditResp, error) {
	body := `{
		"AuditName":"C2C",
		"ContentType":"Text",
		"Content":"%s"
	}`

	// 每次都生成sig 也不要紧
	sig, _ := GetUserSig(appId, key, identifier)
	log.Println("UserSig", sig)

	url := fmt.Sprintf(auditAPI, appId, identifier, sig, rand.Int31())
	log.Println("API", url)

	req, _ := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf(body, text)))
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("IMAuditImage post api err", err.Error())
		return nil, err
	}
	rspV := new(AuditResp)
	err = json.NewDecoder(resp.Body).Decode(rspV)
	if err != nil {
		log.Println("IMAuditImage json.Decode err", err.Error())
		return nil, err
	}
	return rspV, nil
}
