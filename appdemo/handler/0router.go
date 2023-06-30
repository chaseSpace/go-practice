package handler

import (
	"github.com/gin-gonic/gin"
	sugar "go_accost/log"
)

func Setup(eng *gin.Engine) {
	r := eng.Group("/")

	r.POST("/user_online", W(DoUserOnline[UserOnlineReq]))
}

type GCtx struct {
	Ctx *gin.Context
}

type handler[T any] func(ctx *GCtx, r *T) (any, error)

func W[T any](f handler[T]) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		req := new(T)
		err := ctx.Bind(req)
		if err != nil {
			sugar.Errorf("parseReq [%T] %v", req, err)
			return
		}
		cc := &GCtx{}

		var res *Rsp
		rsp, err := f(cc, req)
		if err != nil {
			res = &Rsp{
				Code: 500,
				Msg:  "ERR: " + err.Error(),
				Data: nil,
			}
		} else {
			if rsp == nil {
				rsp = struct{}{}
			}
			res = &Rsp{
				Code: 200,
				Msg:  "OK",
				Data: rsp,
			}
		}
		ctx.JSON(200, res)
		ctx.Next()
	}
}
