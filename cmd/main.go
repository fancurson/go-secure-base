package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	HttpReadTimeOut  = 5 * time.Second
	HttpWriteTimeout = 10 * time.Second
	HttpIdleTimeout  = 120 * time.Second
)

func greet(w http.ResponseWriter, r *http.Request) {
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
		ReadTimeout:  HttpReadTimeOut,  // Ограничение времени на чтение запроса
		WriteTimeout: HttpWriteTimeout, // Ограничение времени на запись ответа
		IdleTimeout:  HttpIdleTimeout,
	}

	log.Println("Server starting on :8080...")

	// Исправляем G104: явно проверяем ошибку
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Critical server failure: %v", err)
	}
}
