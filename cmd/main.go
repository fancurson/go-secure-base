package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	// Указываем Content-Type явно, чтобы избежать атак типа MIME-sniffing
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "GoSecure-Base v1.0 Alive at: %s", time.Now().Format(time.RFC3339))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greet)

	// Настраиваем сервер с защитными таймаутами
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // Защита от медленного чтения (Slowloris)
		WriteTimeout: 10 * time.Second, // Защита от медленной записи
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Server is running on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
