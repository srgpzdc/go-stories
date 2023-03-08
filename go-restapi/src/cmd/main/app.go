package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"restapi/internal/config"
	"restapi/internal/users"
	"restapi/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("create user handler")
	handler := users.NewHandler(logger)
	handler.Register(router)

	logger.Info("startup")
	startup(router, cfg)
}

func startup(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listerErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("listen unix socket")
		listener, listerErr = net.Listen("unix", socketPath)
	} else {
		address := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
		logger.Infof("listen tcp on %s", address)
		listener, listerErr = net.Listen("tcp", address)
	}

	if listerErr != nil {
		logger.Fatal(listerErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
