package main

import (
	"errors"
	"fmt"
	_ "habitgobackend/cmd/api/resource/common/error"
	"habitgobackend/cmd/api/router"
	"habitgobackend/cmd/config"
	"log"
	"net/http"
)

// @title			Atomic Habits Go Backend API
// @version		0.1
// @description	This is the GO backend CRUD REST API for Atomic Habits.
// @basePath		/v1
func main() {
	habitsConfig := config.New()
	routerConfig := router.New()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", habitsConfig.Server.Port),
		Handler:      routerConfig,
		ReadTimeout:  habitsConfig.Server.TimeoutRead,
		WriteTimeout: habitsConfig.Server.TimeoutWrite,
		IdleTimeout:  habitsConfig.Server.TimeoutIdle,
	}

	log.Println("Starting server " + server.Addr)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Server startup failed")
	}
}
