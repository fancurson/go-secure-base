package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	// Безопасность: предотвращаем sniffing атак, явно указывая тип
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "GoSecure-Base: Authorized Access at %s", time.Now().Format(time.RFC3339))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greet)

	// Настраиваем сервер с защитой от DoS (Slowloris)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // Ограничение времени на чтение запроса
		WriteTimeout: 10 * time.Second, // Ограничение времени на запись ответа
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server starting on :8080...")

	// Исправляем G104: явно проверяем ошибку
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Critical server failure: %v", err)
	}
}
