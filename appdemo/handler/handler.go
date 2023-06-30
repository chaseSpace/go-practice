package handler

import sugar "go_accost/log"

/*
Handler func命名以Do开头
*/

func DoUserOnline[T UserOnlineReq](c *GCtx, in *UserOnlineReq) (any, error) {
	sugar.Info("111", in)

	return &UserOnlineRes{Uid: in.Uid}, nil
}
