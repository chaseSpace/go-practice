// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package hello

import (
	"context"
)

type IHelloV1 interface {
	Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error)
	Hello2(ctx context.Context, req *v1.Hello2Req) (res *v1.Hello2Res, err error)
	Hello3(ctx context.Context, req *v1.Hello3Req) (res *v1.Hello3Res, err error)
}
