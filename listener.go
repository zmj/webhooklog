package main

import (
	"fmt"
	"log"
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
	resp.Header().Add("Content-Type", "text/plain")
	resp.Header().Add("Transfer-Encoding", "chunked")
	resp.Header().Add("X-Content-Type-Options", "nosniff")
	resp.WriteHeader(http.StatusOK)
	resp.(http.Flusher).Flush()

	ls := newListener(logID, resp)
	srv.router.add(ls)
	defer srv.router.remove(ls)
	ls.listen()
}

type listener struct {
	logID string
	resp  http.ResponseWriter
	msgs  chan []byte
}

func (ls *listener) listen() {
	for msg := range ls.msgs {
		_, err := ls.resp.Write(msg)
		if err != nil {
			return
		}
		ls.resp.(http.Flusher).Flush()
	}
}

func (ls *listener) log(msg []byte) {
	select {
	case ls.msgs <- msg:
	default:
		log.Println("queue full, dropped message")
	}
}

func newListener(logID string, resp http.ResponseWriter) *listener {
	return &listener{
		logID: logID,
		resp:  resp,
		msgs:  make(chan []byte, 16),
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
