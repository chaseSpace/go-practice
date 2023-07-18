package main

import (
	"fmt"
	"reflect"
)

type BaseRsp struct {
	Code int    `json:"code"` // 默认200为成功，其他均为异常
	Msg  string `json:"msg"`  // 成功为OK，否则为err信息
}

func NewOkRsp() BaseRsp {
	return BaseRsp{
		Code: 200,
		Msg:  "OK",
	}
}

type GetUserInfoRes struct {
	BaseRsp
	IsOnline  bool  `json:"is_online"`
	OnlineSec int32 `json:"online_sec"` // 在线持续时长，秒
}

func main() {
	res := &GetUserInfoRes{}
	reflect.ValueOf(res).Elem().FieldByName("BaseRsp").Set(reflect.ValueOf(NewOkRsp()))
	fmt.Println(res)
}
