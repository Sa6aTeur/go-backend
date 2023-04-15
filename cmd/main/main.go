package main

import (
	"go-backend/internal/user"
	"go-backend/pkg/logger"
	"log"
	"net"
	"net/http"
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

func start(router *httprouter.Router, logger logger.Logger) {
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("hello"))
	})

	listener, err := net.Listen("tcp", "0.0.0.0:1234")

	if err != nil {
		logger.Error("listener cant started")
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("Server start on 0.0.0.0:1234")
	log.Fatalln(server.Serve(listener))
}
