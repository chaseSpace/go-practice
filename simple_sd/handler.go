package simple_sd

import (
	"encoding/json"
	"go-practice/simple_sd/core"
	"net/http"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	instance := core.ServiceInstance{}

	var err error
	defer func() {
		if err != nil {
			core.Sdlogger.Error("handleRegister: service:%s instance:%s error: %v",
				instance.Service, instance.Addr(), err)
			return
		}
		core.Sdlogger.Info("handleRegister OK, service:%s instance:%s", instance.Service, instance.Addr())
	}()

	err = json.NewDecoder(r.Body).Decode(&instance)
	if err != nil {
		return
	}
	_ = r.Body.Close()

	if err = instance.Check(); err != nil {
		return
	}
	err = core.Sd.Register(instance)
}

type discoveryBody struct {
	Service  string
	LastHash string
	Block    bool
}
