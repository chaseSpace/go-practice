package simple_sd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go-practice/simple_sd/core"
	"net/http"
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

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}
	instanceReq := core.ServiceInstance{}

	var data []core.ServiceInstance
	var err error
	var code = 200
	defer func() {
		rsp := ToJson(newRes(data, code, err))
		w.Write(rsp)
		if err != nil {
			core.Sdlogger.Error("handleRegister: service:%s instanceReq:%s error: %v", instanceReq.Service, instanceReq.Addr(), err)
			return
		}
		core.Sdlogger.Info("handleRegister OK, service:%s instanceReq:%s", instanceReq.Service, instanceReq.Addr())
	}()

	err = json.NewDecoder(r.Body).Decode(&instanceReq)
	if err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	_ = r.Body.Close()

	if err = instanceReq.Check(); err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	err = core.Sd.Register(instanceReq)
	if err == nil {
		data, err = core.Sd.Discovery(context.TODO(), instanceReq.Service, "", false)
	}
}

type deregisterBody struct {
	Service string
	Host    string
	Port    int
}

func handleDeregister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}
	req := new(deregisterBody)

	var err error
	var code = 200
	defer func() {
		rsp := ToJson(newRes(nil, code, err))
		w.Write(rsp)
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

type discoveryBody struct {
	Service  string
	LastHash string
	Block    bool
}

func handleDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}

	body := new(discoveryBody)

	var err error
	var code = 200
	var instances []core.ServiceInstance

	defer func() {
		rsp := ToJson(newRes(instances, code, err))
		w.Write(rsp)
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
	instances, err = core.Sd.Discovery(context.TODO(), body.Service, body.LastHash, body.Block)
}
