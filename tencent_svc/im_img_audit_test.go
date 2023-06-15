package tencent_svc

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

type imSdkConf struct {
	AppId      int
	Key        string
	Identifier string
}

func TestIMAuditImage(t *testing.T) {
	b, err := os.ReadFile("../_ignore/im_audit.json")
	if err != nil {
		log.Println("read file err", err.Error())
		return
	}
	sdkConf := new(imSdkConf)
	err = json.Unmarshal(b, sdkConf)
	if err != nil {
		log.Println("Unmarshal err", err.Error())
		return
	}
	u := "https://img.xsnvshen.co/thumb_205x308/album/20946/23797/cover.jpg"
	resp, err := IMAuditImage(u, sdkConf.AppId, sdkConf.Key, sdkConf.Identifier)
	if err != nil {
		log.Printf("Unmarshal err: %+v", err.Error())
		return
	}

	log.Printf("resp: %+v", *resp)
}

func TestIMAuditText(t *testing.T) {
	b, err := os.ReadFile("../_ignore/im_audit.json")
	if err != nil {
		log.Println("read file err", err.Error())
		return
	}
	sdkConf := new(imSdkConf)
	err = json.Unmarshal(b, sdkConf)
	if err != nil {
		log.Println("Unmarshal err", err.Error())
		return
	}
	resp, err := IMAuditText("习大大", sdkConf.AppId, sdkConf.Key, sdkConf.Identifier)
	if err != nil {
		log.Printf("Unmarshal err: %+v", err.Error())
		return
	}

	log.Printf("resp: %+v", *resp)
}
