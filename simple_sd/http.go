package simple_sd

import (
	"fmt"
	"go-practice/simple_sd/core"
	"net/http"
)

type HTTPServer struct {
}

func NewSimpleSdHTTPServer() *HTTPServer {
	core.InitLogger(core.LogLevelDebug)
	return &HTTPServer{}
}

func (s *HTTPServer) Run(port int) {
	http.HandleFunc("/service/register", handleRegister)
	http.HandleFunc("/service/deregister", handleRegister)
	http.HandleFunc("/service/discovery", handleRegister)

	core.Sdlogger.Info("SimpleSdHTTPServer is running on http://localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
