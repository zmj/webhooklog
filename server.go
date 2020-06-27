package main

import (
	"fmt"
	"net/http"
)

type server struct {
}

func (srv *server) serve(port int) error {
	http.HandleFunc("/", srv.defaultHandler)
	http.HandleFunc("/log/", srv.logHandler)
	addr := fmt.Sprintf(":%v", port)
	return http.ListenAndServe(addr, nil)
}

func (srv *server) defaultHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("hello world"))
}
