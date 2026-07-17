package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	VERSION = "v1"
)

func newProxy(target string) *httputil.ReverseProxy {
	u, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(u)
}

func main() {
	users := newProxy("http://localhost:8081")
	orders := newProxy("http://localhost:8082")

	mux := http.NewServeMux()
	mux.Handle("/users/"+VERSION, users)
	mux.Handle("/auth/"+VERSION, orders)

	log.Println("gateway on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
