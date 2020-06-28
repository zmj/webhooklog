package main

import (
	"fmt"
	"net/http"
	"net/url"
)

type server struct {
	router *router
}

func (srv *server) serve(port int) error {
	http.HandleFunc("/log/", srv.logHandler)
	http.HandleFunc("/", srv.defaultHandler)
	http.HandleFunc("/ms/", srv.msGraphHandler)
	addr := fmt.Sprintf(":%v", port)
	return http.ListenAndServe(addr, nil)
}

func (srv *server) defaultHandler(resp http.ResponseWriter, req *http.Request) {
	logID, ok := logIDFromQuery(req.URL)
	if !ok {
		http.Error(resp, "Not Found", http.StatusNotFound)
		return
	}
	srv.router.log(logID, req)

	var status int
	switch req.Method {
	case http.MethodGet:
		status = http.StatusOK
	case http.MethodPost:
		status = http.StatusAccepted
	default:
		status = http.StatusOK
	}
	resp.WriteHeader(status)
}

func logIDFromQuery(url *url.URL) (string, bool) {
	logID := url.Query().Get("log")
	return logID, logID != ""
}

func newServer() *server {
	return &server{router: newRouter()}
}
