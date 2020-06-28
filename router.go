package main

import "net/http"

type router struct{}

func (r *router) log(logID string, req *http.Request) {

}

func newRouter() *router {
	return &router{}
}
