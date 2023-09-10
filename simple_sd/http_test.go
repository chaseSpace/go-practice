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
	"time"
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
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(registerReq)).Error())),
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
			wantRes: newRes([]*core.ServiceInstance{
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
			t.Fatalf(tt.name+"--handler res got：%s, not want:%+v", b, tt.wantRes)
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

// 对比耗时的时候必须考虑时间误差
func isMillsTimeCostEqual(got, want int64) bool {
	return got >= want && got <= want+core.DiscoveryInterval.Milliseconds()
}

func TestHTTPServer_handleDiscovery(t *testing.T) {
	type XX struct {
		name              string
		req               interface{}
		regAfter          time.Duration
		regInstances      []*core.ServiceInstance
		regInstancesAfter []*core.ServiceInstance
		wantRes           *HttpRes
	}

	__inst1 := []*core.ServiceInstance{{
		Service: "go-user",
		Host:    "localhost",
		Port:    8080,
	}}
	__inst2 := []*core.ServiceInstance{{
		Service: "go-user",
		Host:    "localhost",
		Port:    8081,
	}}

	__oneInstHash := core.CalInstanceHash(__inst1)
	__twoInstHash := core.CalInstanceHash(append(__inst1, __inst2...))
	tests := []XX{
		{
			name:    "T-无效body",
			req:     []int{1},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, json.Unmarshal([]byte(`[1]`), new(discoveryReq)).Error())),
		},
		{
			name: "T-空service",
			req: &discoveryReq{
				Service: "",
			},
			wantRes: newRes(nil, 400, errors.Wrap(ErrParams, "need service")),
		},
		{
			name: "T-设置Service",
			req: &discoveryReq{
				Service: "go-user",
			},
			regInstances: __inst1,
			wantRes: newRes(discoveryRsp{
				Instances: __inst1,
				Hash:      core.CalInstanceHash(__inst1),
			}, 200, nil),
		},
		{
			name: "T-设置Service、空Hash、WaitMaxMs=10, 无实例（应耗尽时间再返回）",
			req: &discoveryReq{
				Service:   "go-user",
				WaitMaxMs: 10,
			},
			wantRes: newRes(discoveryRsp{}, 200, nil),
		},
		{
			name: "T-设置Service、空Hash、WaitMaxMs=10, 有实例（应立即返回）",
			req: &discoveryReq{
				Service:   "go-user",
				WaitMaxMs: 10,
			},
			regInstances: __inst1,
			wantRes: newRes(discoveryRsp{
				Instances: __inst1,
				Hash:      core.CalInstanceHash(__inst1),
			}, 200, nil),
		},
		{
			name: "T-设置Service、与实例一致的Hash、WaitMaxMs=10, 有实例（应耗尽时间再返回）",
			req: &discoveryReq{
				Service:   "go-user",
				WaitMaxMs: 10,
				LastHash:  core.CalInstanceHash(__inst1),
			},
			regInstances: __inst1,
			wantRes: newRes(discoveryRsp{
				Instances: __inst1,
				Hash:      core.CalInstanceHash(__inst1),
			}, 200, nil),
		},
		{
			name: "T-设置Service、空Hash、WaitMaxMs=15, 10ms后注册1个实例（应耗尽时间再返回1个实例）",
			req: &discoveryReq{
				Service:   "go-user",
				WaitMaxMs: 20, // 需要稍微大于 regAfter
			},
			regAfter:          time.Millisecond * 10,
			regInstancesAfter: __inst1,
			wantRes: newRes(discoveryRsp{
				Instances: __inst1,
				Hash:      __oneInstHash,
			}, 200, nil),
		},
		{
			name: "T-设置Service、与1个实例一致的Hash、WaitMaxMs=1000, 先注册1个实例，1000ms后再注册另一个实例（应耗尽时间再返回2个实例）",
			req: &discoveryReq{
				Service:   "go-user",
				WaitMaxMs: 1010, // 需要稍微大于 regAfter
				LastHash:  __oneInstHash,
			},
			regInstances:      __inst1,
			regAfter:          time.Millisecond * 1000,
			regInstancesAfter: __inst2,
			wantRes: newRes(discoveryRsp{
				Instances: append(__inst1, __inst2...),
				Hash:      __twoInstHash,
			}, 200, nil),
		},
	}

	for _, tt := range tests {
		for _, inst := range tt.regInstances {
			err := core.Sd.Register(*inst)
			if err != nil {
				t.Fatalf(tt.name+"--- Register failed %v", err)
			}
		}
		go func(instances []*core.ServiceInstance) {
			time.Sleep(tt.regAfter)
			for _, inst := range instances {
				err := core.Sd.Register(*inst)
				if err != nil {
					t.Fatalf(tt.name+"--- Register failed %v", err)
				}
			}
		}(tt.regInstancesAfter)

		req, err := http.NewRequest("POST", "/service/discovery", bytes.NewReader(ToJson(tt.req)))
		if err != nil {
			t.Fatal(tt.name, err)
		}

		rr := httptest.NewRecorder()

		st := time.Now()
		handler := http.HandlerFunc(handleDiscovery)
		handler.ServeHTTP(rr, req)
		cost := time.Since(st).Milliseconds()

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf(tt.name+"--handler unexpected statuscode: got %v, want %v", status, http.StatusOK)
		}

		body, ok := tt.req.(*discoveryReq)
		if ok && body.WaitMaxMs > 0 && !isMillsTimeCostEqual(cost, tt.regAfter.Milliseconds()) {
			t.Fatalf(tt.name+"--handler unexpected waitMs: got %v, want %v", cost, tt.regAfter.Milliseconds())
		}

		// default time cost limited to 2ms
		if ok && body.WaitMaxMs == 0 && cost > 2 {
			println(222)
			t.Fatalf(tt.name+"--handler unexpected waitMs: got %v, want %v", cost, 2)
		}

		b := rr.Body.Bytes()
		if !bytes.Equal(ToJson(tt.wantRes), b) {
			t.Fatalf(tt.name+"--handler res got:\n%s\n--- not want:\n%s", b, ToJson(tt.wantRes))
		}

		for _, inst := range append(__inst1, __inst2...) {
			_ = core.Sd.Deregister(inst.Service, inst.Addr())
		}

		//fmt.Printf(tt.name+"--handler want:%+v\n", tt.wantRes)
	}
}

func doRegister(insts []*core.ServiceInstance) []byte {
	var res []byte
	for _, inst := range insts {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleRegister)
		req, _ := http.NewRequest("POST", "/service/register", bytes.NewReader(ToJson(registerReq{*inst})))
		handler.ServeHTTP(rr, req)
		res = rr.Body.Bytes()
	}
	return res
}

func doDiscovery(service string) []byte {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleDiscovery)
	req, _ := http.NewRequest("POST", "/service/discovery", bytes.NewReader(ToJson(&discoveryReq{
		Service: service,
	})))
	handler.ServeHTTP(rr, req)
	return rr.Body.Bytes()
}

func mustHaveInstance(t *testing.T, insts []*core.ServiceInstance, hash string) {
	var wantRes = &discoveryRsp{
		Instances: insts,
		Hash:      hash,
	}
	resBytes := doDiscovery("go-user")

	if !bytes.Equal(resBytes, ToJson(newRes(wantRes, 200, nil))) {
		t.Fatalf("Discovery result failed, got res:%s", resBytes)
	}
}

func mustNoInstance(t *testing.T) {
	var wantRes = &discoveryRsp{}
	resBytes := doDiscovery("go-user")

	if !bytes.Equal(resBytes, ToJson(newRes(wantRes, 200, nil))) {
		t.Fatalf("Discovery result failed, got res:%s", resBytes)
	}
}

func TestHTTPServer_handleDiscoveryWithHealthCheckFail(t *testing.T) {
	__inst1 := []*core.ServiceInstance{{
		Service: "go-user",
		Host:    "localhost",
		Port:    8080,
	}}
	//__inst2 := []*core.ServiceInstance{{
	//	Service: "go-user",
	//	Host:    "localhost",
	//	Port:    8081,
	//}}

	__oneInstHash := core.CalInstanceHash(__inst1)
	//__twoInstHash := core.CalInstanceHash(append(__inst1, __inst2...))

	resBytes := doRegister(__inst1)
	if !bytes.Equal(resBytes, ToJson(newRes(__inst1, 200, nil))) {
		t.Fatalf("doRegister failed: %s", resBytes)
	}

	// 在sd进行健康检测期间，实例应该一直存在
	serverAliveSec := core.HealthCheckInterval.Seconds() * core.HealthCheckMaxFails
	for i := 0; i < int(serverAliveSec); i++ {
		println("--- No.1 sleep serverAliveSec", i+1)
		mustHaveInstance(t, __inst1, __oneInstHash)
		time.Sleep(time.Second)
	}

	println("No.1 sleep end...")
	time.Sleep(time.Second) // 给一个合理的间隙时间
	mustNoInstance(t)
}

func TestHTTPServer_handleDiscoveryWithHealthCheckPass(t *testing.T) {
	__inst1 := []*core.ServiceInstance{{
		Service: "go-user",
		Host:    "localhost",
		Port:    8080,
	}}
	//__inst2 := []*core.ServiceInstance{{
	//	Service: "go-user",
	//	Host:    "localhost",
	//	Port:    8081,
	//}}

	__oneInstHash := core.CalInstanceHash(__inst1)
	//__twoInstHash := core.CalInstanceHash(append(__inst1, __inst2...))

	//
	resBytes := doRegister(__inst1)
	// -- 返回实例是按注册时间排序的，所以这里顺序固定
	if !bytes.Equal(resBytes, ToJson(newRes(__inst1, 200, nil))) {
		t.Fatalf("doRegister failed: %s", resBytes)
	}

	// -- 同时启动实例 以通过健康检测
	__server := http.Server{Addr: __inst1[0].Addr()}
	serverAliveSec := core.HealthCheckInterval.Seconds() * (core.HealthCheckMaxFails + 1)
	go func() {
		err := __server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}
		println("No.2 server closed")
	}()

	time.AfterFunc(time.Duration(serverAliveSec)*time.Second, func() {
		__server.Close()
	})

	// 以更长的时间来持续检测实例是否注册
	for i := 0; i < int(serverAliveSec); i++ {
		println("--- No.2 sleep serverAliveSec", i+1)
		mustHaveInstance(t, __inst1, __oneInstHash)
		time.Sleep(time.Second)
	}

	println("No.2 sleep end...")
	time.Sleep(time.Second)
	mustNoInstance(t)
}
