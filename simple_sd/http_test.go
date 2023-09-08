package simple_sd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go-practice/simple_sd/core"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHTTPServer_handleRegister(t *testing.T) {
	type XX struct {
		name    string
		req     interface{}
		wantRes *HttpRes
	}
	tests := []XX{
		{
			name:    "T-无效body",
			req:     []int{1},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(core.ServiceInstance)).Error())),
		},
		{
			name: "T-空service",
			req: core.ServiceInstance{
				Service: "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, core.ServiceInstance{}.Check().Error())),
		},
		{
			name: "T-空Host",
			req: core.ServiceInstance{
				Service: "go-user",
				Host:    "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, core.ServiceInstance{Service: "go-user", Host: ""}.Check().Error())),
		},
		{
			name: "T-空Port",
			req: core.ServiceInstance{
				Service: "go-user",
				Host:    "localhost",
				Port:    0,
			},
			wantRes: newRes(nil, 400,
				errors.Wrap(ErrParams, core.ServiceInstance{Service: "go-user", Host: "localhost", Port: 0}.Check().Error())),
		},
		{
			name: "T-有效实例",
			req: core.ServiceInstance{
				Service: "go-user",
				Host:    "localhost",
				Port:    8080,
			},
			wantRes: newRes([]core.ServiceInstance{
				{
					Service: "go-user",
					Host:    "localhost",
					Port:    8080,
				},
			}, 200, nil),
		},
	}

	for _, tt := range tests {
		// 创建一个模拟的HTTP请求
		req, err := http.NewRequest("POST", "/service/register", bytes.NewReader(ToJson(tt.req)))
		if err != nil {
			t.Fatal(tt.name, err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handleRegister)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf(tt.name+"--handler unexpected statuscode: got %v, want %v", status, http.StatusOK)
		}

		b := rr.Body.Bytes()
		if !bytes.Equal(ToJson(tt.wantRes), b) {
			t.Fatalf(tt.name+"--handler res got：%+v, not want:%+v", b, tt.wantRes)
		}
		//fmt.Printf(tt.name+"--handler want:%+v\n", tt.wantRes)
	}
}

func checkItem(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("---NotEqual a:%+v b:%+v", a, b)
	}
}

func TestHTTPServer_handleDeregister(t *testing.T) {
	type XX struct {
		name     string
		req      interface{}
		regFirst bool
		wantRes  *HttpRes
	}
	tests := []XX{
		{
			name:    "T-无效body",
			req:     []int{1},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(deregisterBody)).Error())),
		},
		{
			name: "T-空service",
			req: deregisterBody{
				Service: "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, "provide a valid service, host, port")),
		},
		{
			name: "T-空Host",
			req: deregisterBody{
				Service: "go-user",
				Port:    8000,
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, "provide a valid service, host, port")),
		},
		{
			name: "T-空Port",
			req: deregisterBody{
				Service: "go-user",
				Host:    "localhost",
				Port:    0,
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, "provide a valid service, host, port")),
		},
		{
			name: "T-注销未注册的实例",
			req: deregisterBody{
				Service: "go-user",
				Host:    "localhost",
				Port:    8000,
			},
			wantRes: newRes(nil, 400, errors.Wrap(core.ErrInstanceNotRegistered,
				fmt.Sprintf("service: %s, addr: %s", "go-user", "localhost:8000"))),
		},
		{
			name: "T-正常注销实例",
			req: deregisterBody{
				Service: "go-user",
				Host:    "localhost",
				Port:    8000,
			},
			regFirst: true,
			wantRes:  newRes(nil, 200, nil),
		},
	}

	for _, tt := range tests {
		if tt.regFirst {
			inst := tt.req.(deregisterBody)
			err := core.Sd.Register(core.ServiceInstance{
				Service:  inst.Service,
				IsUDP:    false,
				Host:     inst.Host,
				Port:     inst.Port,
				Metadata: nil,
			})
			if err != nil {
				t.Fatalf(tt.name+"--- Register failed %v", err)
			}
		}
		req, err := http.NewRequest("POST", "/service/deregister", bytes.NewReader(ToJson(tt.req)))
		if err != nil {
			t.Fatal(tt.name, err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handleDeregister)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf(tt.name+"--handler unexpected statuscode: got %v, want %v", status, http.StatusOK)
		}

		b := rr.Body.Bytes()
		if !bytes.Equal(ToJson(tt.wantRes), b) {
			t.Fatalf(tt.name+"--handler res got:\n%s\n--- not want:\n%s", b, ToJson(tt.wantRes))
		}
		//fmt.Printf(tt.name+"--handler want:%+v\n", tt.wantRes)
	}
}

func TestHTTPServer_handleDiscovery(t *testing.T) {
	type XX struct {
		name         string
		req          interface{}
		regInstances []core.ServiceInstance
		wantRes      *HttpRes
	}
	tests := []XX{
		{
			name:    "T-无效body",
			req:     []int{1},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(discoveryBody)).Error())),
		},
		{
			name: "T-空service",
			req: discoveryBody{
				Service: "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, "need service")),
		},
	}

	for _, tt := range tests {
		for _, inst := range tt.regInstances {
			err := core.Sd.Register(core.ServiceInstance{
				Service:  inst.Service,
				IsUDP:    false,
				Host:     inst.Host,
				Port:     inst.Port,
				Metadata: nil,
			})
			if err != nil {
				t.Fatalf(tt.name+"--- Register failed %v", err)
			}
		}
		req, err := http.NewRequest("POST", "/service/discovery", bytes.NewReader(ToJson(tt.req)))
		if err != nil {
			t.Fatal(tt.name, err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(handleDiscovery)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf(tt.name+"--handler unexpected statuscode: got %v, want %v", status, http.StatusOK)
		}

		b := rr.Body.Bytes()
		if !bytes.Equal(ToJson(tt.wantRes), b) {
			t.Fatalf(tt.name+"--handler res got:\n%s\n--- not want:\n%s", b, ToJson(tt.wantRes))
		}
		//fmt.Printf(tt.name+"--handler want:%+v\n", tt.wantRes)
	}
}
