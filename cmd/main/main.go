package main

import (
	"fmt"
	"go-backend/internal/config"
	"go-backend/internal/user"
	"go-backend/pkg/logger"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger.Init()
	logger := logger.GetLogger()
	router := httprouter.New()

	userHandler := user.NewHandler()
	userHandler.Register(router)

	start(router, logger)
}

func start(router *httprouter.Router, logger *logger.Logger) {
	var listenerErr error
	var listener net.Listener
	cfg := config.GetConfig()

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatalln(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("Create socket. Socket path: %s", socketPath)
		logger.Info("Create unix socket")
		listener, listenerErr = net.Listen("unix", socketPath)
	} else {
		logger.Info("Create tcp listener")
		listener, listenerErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	}

	if listenerErr != nil {
		logger.Fatalln(listenerErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("Server start on %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	log.Fatalln(server.Serve(listener))
}
