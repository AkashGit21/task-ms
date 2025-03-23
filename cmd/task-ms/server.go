package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	taskV1 "github.com/AkashGit21/task-ms/api/task/v1"
	"github.com/AkashGit21/task-ms/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func NewTaskV1Server() (*http.Server, error) {
	taskV1API, err := taskV1.New()
	if err != nil {
		return nil, err
	}

	return newServer(taskV1API)
}

func NewAuthnServer() (*http.Server, error) {
	authnAPI, err := taskV1.New()
	if err != nil {
		return nil, err
	}

	return newServer(authnAPI)
}

func StartServer(srv *http.Server) {
	utils.InfoLog("Starting Server...")

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down the server gracefully...")
		// Received an interrupt signal, shut down once the current requests are finished
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Println("HTTP server Shutdown: ", err)
			return
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Println("HTTP server ListenAndServe: ", err)
		return
	}

	<-idleConnsClosed
}

func newServer(api *mux.Router) (*http.Server, error) {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	srvHost := utils.GetEnvValue("APP_HOST", "localhost")
	srvPort := utils.GetEnvValue("APP_PORT", "8081")
	srvAddress := fmt.Sprintf("%s:%v", srvHost, srvPort)
	log.Println("Configuring Server at address", srvAddress)

	srv := http.Server{
		Addr:    srvAddress,
		Handler: api,
		// Read will Timeout after 2s if anything goes wrong.
		ReadTimeout: time.Duration(2 * time.Second),
		// Write will Timeout after 5s
		WriteTimeout: time.Duration(5 * time.Second),
	}

	return &srv, nil
}
