package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (srv *server) logHandler(resp http.ResponseWriter, req *http.Request) {
	logID, err := logIDFromPath(req.URL)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Write([]byte(logID))
}

func logIDFromPath(url *url.URL) (string, error) {
	path := strings.Trim(url.Path, "/")
	segments := strings.Split(path, "/")
	if len(segments) != 2 {
		return "", fmt.Errorf("invalid path: expected '/log/<logID>'")
	}
	return segments[1], nil
}
