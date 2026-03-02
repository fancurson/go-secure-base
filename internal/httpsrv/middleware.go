// Package httpsrv provides HTTP server implementation and middlewares.
package httpsrv

import "net/http"

// SecurityHeaders injected mandatory security headers into the response.
// This is a defensive layer protecting against XSS, Clickjacking, and MIME-sniffing.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. HSTS: Принудительный HTTPS на 1 год (включая поддомены)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 2. CSP: Запрещаем всё (идеально для чистого JSON API)
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; sandbox;")

		// 3. X-Frame-Options: Защита от кликджекинга
		w.Header().Set("X-Frame-Options", "DENY")

		// 4. X-Content-Type-Options: Запрет угадывания типов (MIME-sniffing)
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// 5. Referrer-Policy: Минимальная утечка данных при переходе по ссылкам
		w.Header().Set("Referrer-Policy", "no-referrer")

		// Передаем управление следующему обработчику в цепочке
		next.ServeHTTP(w, r)
	})
}
