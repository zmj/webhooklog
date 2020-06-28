package main

import (
	"fmt"
	"log"
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
	resp.Header().Add("Cache-Control", "no-store")
	resp.Header().Add("Cache-Control", "no-transform")
	resp.Write([]byte{newline})
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
	var err error
	for err == nil {
		select {
		case msg := <-ls.msgs:
			_, err = ls.resp.Write([]byte{newline})
			_, err = ls.resp.Write(msg)
			_, err = ls.resp.Write([]byte{newline})
		case <-time.After(10 * time.Second):
			_, err = ls.resp.Write([]byte{period})
		}
		ls.resp.(http.Flusher).Flush()
	}
}

const newline = byte('\n')
const period = byte('.')

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
