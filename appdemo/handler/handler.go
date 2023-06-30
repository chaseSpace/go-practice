package handler

import sugar "go_accost/log"

/*
Handler func命名以Do开头
*/

func DoUserOnline(c *GCtx) (any, error) {
	req := c.Req.(*UserOnlineReq)
	sugar.Info("111", req)

	return &UserOnlineRes{Uid: req.Uid}, nil
}

// TODO 用户注册通知
