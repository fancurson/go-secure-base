// Package httpsrv provides the HTTP server initialization and routing.
package httpsrv

import (
	"context"
	"net"
	"net/http"

	"github.com/fancurson/go-secure-base/internal/config"
)

// NewServer returns a pre-configured http.Server with hardened settings.
func NewServer(baseCtx context.Context, cfg *config.Config, handler http.Handler) *http.Server {
	SecureHandler := SecurityHeaders(handler)

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           SecureHandler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,

		BaseContext: func(_ net.Listener) context.Context {
			return baseCtx
		},
	}
}
