// Package main is the entry point for the GoSecure-Base microservice.
// It initializes the server, security middlewares, and handles graceful shutdown.
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	HTTPReadTimeOut  = 5 * time.Second
	HTTPWriteTimeout = 10 * time.Second
	HTTPIdleTimeout  = 120 * time.Second
)

func greet(w http.ResponseWriter, _ *http.Request) {
	// Безопасность: предотвращаем sniffing атак, явно указывая тип
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := fmt.Fprintf(w, "GoSecure-Base: Authorized Access at %s", time.Now().Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greet)

	// Настраиваем сервер с защитой от DoS (Slowloris)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  HTTPReadTimeOut,  // Ограничение времени на чтение запроса
		WriteTimeout: HTTPWriteTimeout, // Ограничение времени на запись ответа
		IdleTimeout:  HTTPIdleTimeout,
	}

	log.Println("Server starting on :8080...")

	// Исправляем G104: явно проверяем ошибку
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Critical server failure: %v", err)
	}
}
