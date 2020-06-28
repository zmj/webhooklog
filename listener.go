package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (srv *server) logHandler(resp http.ResponseWriter, req *http.Request) {
	logID, err := logIDFromPath(req.URL)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Header().Add("Content-Type", "text/plain")
	resp.Header().Add("Transfer-Encoding", "chunked")
	resp.Header().Add("X-Content-Type-Options", "nosniff")
	for err == nil {
		_, err = resp.Write([]byte(logID))
		resp.(http.Flusher).Flush()
		<-time.After(time.Second)
	}
}

func logIDFromPath(url *url.URL) (string, error) {
	path := strings.Trim(url.Path, "/")
	segments := strings.Split(path, "/")
	if len(segments) != 2 {
		return "", fmt.Errorf("invalid path: expected '/log/<logID>'")
	}
	return segments[1], nil
}
