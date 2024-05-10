package hello

import (
	"context"
	v1 "webapp/api/hello/v1"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) Hello2(ctx context.Context, req *v1.Hello2Req) (res *v1.Hello2Res, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) Hello3(ctx context.Context, req *v1.Hello3Req) (res *v1.Hello3Res, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
