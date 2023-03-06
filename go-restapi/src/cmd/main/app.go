package main

import (
	"net"
	"net/http"
	"restapi/internal/users"
	"restapi/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	logger.Info("create user handler")
	handler := users.NewHandler(logger)
	handler.Register(router)

	logger.Info("startup")
	startup(router)
}

func startup(router *httprouter.Router) {
	logger := logging.GetLogger()
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
