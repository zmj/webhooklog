package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
)

type router struct {
	mu        sync.Mutex
	listeners map[string][]*listener
}

func (r *router) log(logID string, req *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lss := r.listeners[logID]
	if len(lss) == 0 {
		return
	}
	msg, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
		return
	}
	for _, ls := range lss {
		ls.log(msg)
	}
}

func (r *router) add(ls *listener) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lss := r.listeners[ls.logID]
	r.listeners[ls.logID] = append(lss, ls)
}

func (r *router) remove(ls *listener) {
	r.mu.Lock()
	defer r.mu.Unlock()
	lss := r.listeners[ls.logID]
	for i, otherLs := range lss {
		if ls == otherLs {
			last := len(lss) - 1
			lss[last], lss[i] = lss[i], lss[last]
			break
		}
	}
	lss = lss[:len(lss)-1]
	if len(lss) == 0 {
		delete(r.listeners, ls.logID)
	} else {
		r.listeners[ls.logID] = lss
	}
}

func newRouter() *router {
	return &router{
		listeners: make(map[string][]*listener),
	}
}
