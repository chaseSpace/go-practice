package tencent_svc

import (
	"encoding/json"
	"fmt"
	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

const auditImageAPI = "https://console.tim.qq.com/v4/im_msg_audit/content_moderation?sdkappid=%d&identifier=%s&usersig=%s&random=%d&contenttype=json"

type AuditResp struct {
	ActionStatus string   `json:"ActionStatus"`
	ErrorCode    int      `json:"ErrorCode"`
	ErrorInfo    string   `json:"ErrorInfo"`
	RequestID    string   `json:"RequestId"`
	Result       string   `json:"Result"`
	Score        int      `json:"Score"`
	Label        string   `json:"Label"`
	Keywords     []string `json:"Keywords"`
}

func GetUserSig(appId int, key, uid string) (string, error) {
	userSig, err := tencentyun.GenUserSig(appId, key, uid, 600)
	if err != nil {
		log.Println("GetUserSig err", err.Error())
		return "", err
	}
	return userSig, nil
}

func IMAuditImage(imgURL string, appId int, key, identifier string) (*AuditResp, error) {
	body := `{
		"AuditName":"C2C",
		"ContentType":"Image",
		"Content":"%s"
	}`

	// 每次都生成sig 也不要紧
	sig, _ := GetUserSig(appId, key, identifier)
	log.Println("UserSig", sig)

	url := fmt.Sprintf(auditImageAPI, appId, identifier, sig, rand.Int31())
	log.Println("API", url)

	req, _ := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf(body, imgURL)))
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
