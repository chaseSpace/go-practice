package simple_sd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go-practice/simple_sd/core"
	"net/http"
	"time"
)

type HttpRes struct {
	Code int // 200 OK
	Msg  string
	Data interface{} `json:"Data,omit_empty"`
}

func ToJson(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

var (
	ErrMethod       = errors.New("use HTTP POST")
	ErrParams       = errors.New("invalid params")
	ErrInstanceAddr = errors.New("invalid Addr, provide a addr like localhost:8000 or 192.168.1.1:8000")
)

func newRes(data interface{}, code int, err error) *HttpRes {
	msg := "OK"
	if err != nil {
		msg = err.Error()
	}
	return &HttpRes{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

type registerReq struct {
	core.ServiceInstance
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}
	req := new(registerReq)

	var data []*core.ServiceInstance
	var err error
	var code = 200
	defer func() {
		rsp := ToJson(newRes(data, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			core.Sdlogger.Error("handleRegister: service:%s req:%s error: %v", req.Service, req.Addr(), err)
			return
		}
		core.Sdlogger.Info("handleRegister OK, service:%s req:%s", req.Service, req.Addr())
	}()

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	_ = r.Body.Close()

	if err = req.Check(); err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	err = core.Sd.Register(req.ServiceInstance)
	if err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()
		data, _, err = core.Sd.Discovery(ctx, req.Service, "")
	}
}

type deregisterBody struct {
	Service string
	Host    string
	Port    int
}

func handleDeregister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}
	req := new(deregisterBody)

	var err error
	var code = 200
	defer func() {
		rsp := ToJson(newRes(nil, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			core.Sdlogger.Error("handleDeregister: req:%+v error: %v", req, err)
			return
		}
		core.Sdlogger.Info("handleDeregister OK, req:%+v", req)
	}()

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	_ = r.Body.Close()
	if req.Service == "" || req.Host == "" || req.Port == 0 {
		code = 400
		err = errors.Wrap(ErrParams, "provide a valid service, host, port")
		return
	}
	err = core.Sd.Deregister(req.Service, fmt.Sprintf("%s:%d", req.Host, req.Port))
	if err != nil {
		if errors.Is(err, core.ErrInstanceNotRegistered) {
			code = 400
		}
	}
}

type discoveryReq struct {
	Service   string
	LastHash  string
	WaitMaxMs int64
}
type discoveryRsp struct {
	Instances []*core.ServiceInstance
	Hash      string
}

const MaxDiscoveryTimeout = time.Minute * 5

func handleDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}

	body := new(discoveryReq)

	var err error
	var code = 200

	var response *discoveryRsp

	defer func() {
		rsp := ToJson(newRes(response, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			core.Sdlogger.Error("handleDiscovery: body:%+v error: %v", body, err)
			return
		}
		core.Sdlogger.Info("handleDiscovery OK, body:%+v", body)
	}()

	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	_ = r.Body.Close()
	if body.Service == "" {
		code = 400
		err = errors.Wrap(ErrParams, "need service")
		return
	}
	if body.WaitMaxMs > MaxDiscoveryTimeout.Milliseconds() {
		code = 400
		err = errors.Wrapf(ErrParams, "max wait ms is %s", MaxDiscoveryTimeout)
		return
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*time.Duration(body.WaitMaxMs))
	defer cancel()

	var (
		instances []*core.ServiceInstance
		hash      string
	)
	instances, hash, err = core.Sd.Discovery(ctx, body.Service, body.LastHash)
	if err == nil {
		response = &discoveryRsp{
			Instances: instances,
			Hash:      hash,
		}
	}
}
