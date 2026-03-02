// Package main is the entry point for the GoSecure-Base microservice.
package main

import (
	"fmt"
	"os"

	"github.com/fancurson/go-secure-base/internal/app"
)

func main() {
	// В будущем тут может быть парсинг флагов
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application failed: %v\n", err)
		os.Exit(1)
	}

}
