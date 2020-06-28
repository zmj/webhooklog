package main

import (
	"net/http"
)

func (srv *server) msGraphHandler(resp http.ResponseWriter, req *http.Request) {
	logID, ok := logIDFromQuery(req.URL)
	if !ok {
		http.Error(resp, "Not Found", http.StatusNotFound)
		return
	}
	srv.router.log(logID, req)

	vt := req.URL.Query().Get("validationToken")
	if vt != "" {
		resp.Header().Add("Content-Type", "text/plain")
		resp.Write([]byte(vt))
	} else {
		resp.WriteHeader(http.StatusAccepted)
	}
}
