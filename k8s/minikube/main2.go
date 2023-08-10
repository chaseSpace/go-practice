package main

import (
	"io"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "[v2] Hello, Kubernetes!")
}

func main() {
	http.HandleFunc("/", hello)

	log.Printf("v2 access http://localhost:3000\n")
	panic(http.ListenAndServe(":3000", nil))
}
