package main

import (
	"log"
	"net"
	"net/http"
	"restapi/internal/users"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	log.Println("create router")
	router := httprouter.New()

	log.Println("create user handler")
	handler := users.NewHandler()
	handler.Register(router)

	log.Println("startup")
	startup(router)
}

func startup(router *httprouter.Router) {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
