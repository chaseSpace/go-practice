package api

import (
	"github.com/parnurzeal/gorequest"
	sugar "go_accost/log"
)

var dor *gorequest.SuperAgent
var endpoint = "http://127.0.0.1:8080/"

func Init(mainSvcAddr string) {
	dor = gorequest.New()
	endpoint = mainSvcAddr
}

func logRequest(name string, in, out interface{}, err []error, hint ...string) {
	_hint := "{0000}---"
	if len(hint) > 0 {
		_hint = hint[0] + "---"
	}
	if err != nil {
		sugar.Error(_hint+name+" err", in, "--err", err)
	} else {
		sugar.Debug(_hint+name+" OK", in, "--out", out)
	}
}

// ----------------------

// GetUserInfo 请求用户信息以及在线状态
func GetUserInfo(in *GetUserInfoReq, out *GetUserInfoRes, hint ...string) []error {
	name := "GetUserInfo"
	sugar.Info(name+" 111", in)

	_, _, err := dor.Get(endpoint + name).EndStruct(&out)
	logRequest(name, in, out, err)
	return err
}
