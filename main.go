package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	port, err := choosePort()
	if err != nil {
		log.Println(err)
	}
	srv := &server{
		router: newRouter(),
	}
	err = srv.serve(port)
	if err != nil {
		log.Println(err)
	}
}

func choosePort() (int, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return 8080, nil
	}
	return strconv.Atoi(args[0])
}
