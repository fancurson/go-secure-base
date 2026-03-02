// Package config provides the server configuration parameters and default settings.
package config

import (
	"time"
)

const (
	httpReadTimeout       = 5 * time.Second
	httpWriteTimeout      = 10 * time.Second
	httpIdleTimeout       = 120 * time.Second
	httpReadHeaderTimeout = 10 * time.Second

	addr = ":8080"

	httpMaxHeaderBytes = 1 << 20 // 1MB
)

// Config represents the server configuration parameters.
type Config struct {
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration

	MaxHeaderBytes int

	Addr string
}

// New returns a new Config with default values.
func New() *Config {
	return &Config{
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
		ReadHeaderTimeout: httpReadHeaderTimeout,

		MaxHeaderBytes: httpMaxHeaderBytes,

		Addr: addr,
	}
}
