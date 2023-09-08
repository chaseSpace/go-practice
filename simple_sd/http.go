package simple_sd

import (
	"fmt"
	"go-practice/simple_sd/core"
	"net/http"
)

type HTTPServer struct {
}

func NewSimpleSdHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Run(port int) {
	http.HandleFunc("/service/register", handleRegister)
	http.HandleFunc("/service/deregister", handleDeregister)
	http.HandleFunc("/service/discovery", handleDiscovery)

	core.Sdlogger.Info("SimpleSdHTTPServer is running on http://localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
