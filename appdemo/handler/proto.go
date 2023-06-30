package handler

type Rsp struct {
	Code int         `json:"code"` // 默认200
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type UserOnlineReq struct {
	Uid int64 `json:"uid"`
}
type UserOnlineRes struct {
	Uid int64 `json:"uid"`
}
