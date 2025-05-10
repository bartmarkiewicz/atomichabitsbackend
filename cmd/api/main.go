package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"habitgobackend/cmd/api/config/router"
	"habitgobackend/cmd/api/config/validation"
	_ "habitgobackend/cmd/api/resource/common/error"
	"habitgobackend/cmd/config"
	"log"
	"net/http"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

// @title			Atomic Habits Go Backend API
// @version		0.1
// @description	This is the GO backend CRUD REST API for Atomic Habits.
// @basePath		/v1
func main() {
	habitsConfig := config.New()
	validator := validation.New()

	var logLevel gormlogger.LogLevel
	if habitsConfig.Server.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	dbString := fmt.Sprintf(fmtDBString, habitsConfig.Database.Host, habitsConfig.Database.Username,
		habitsConfig.Database.Password,
		habitsConfig.Database.DatabaseName, habitsConfig.Database.Port)
	database, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		log.Fatal("DB connection start failure")
		return
	}

	routerConfig := router.New(database, validator)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", habitsConfig.Server.Port),
		Handler:      routerConfig,
		ReadTimeout:  habitsConfig.Server.TimeoutRead,
		WriteTimeout: habitsConfig.Server.TimeoutWrite,
		IdleTimeout:  habitsConfig.Server.TimeoutIdle,
	}

	log.Println("Starting server " + server.Addr)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Server startup failed")
	}
}
