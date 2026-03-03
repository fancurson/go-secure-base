// Package logging provides infrastructure components for application logging,
// including security middlewares for sensitive data redaction.
package logging

import (
	"context"
	"log/slog"
	"strings"
)

// sensitiveKeys contains a list of keys that are prohibited from logging in the clear
var sensitiveKeys = map[string]struct{}{
	"password":      {},
	"token":         {},
	"auth":          {},
	"authorization": {},
	"secret":        {},
	"cookie":        {},
	"api_key":       {},
}

// SecurityHandler is a middleware for slog.Handler that automatically
// masks sensitive information in log attributes based on a blacklist.
type SecurityHandler struct {
	Next slog.Handler
}

// Handle processes the log record by masking sensitive attribute values,
// ensuring data protection before the record is passed to the next handler.
func (s SecurityHandler) Handle(ctx context.Context, r slog.Record) error {
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	r.Attrs(func(a slog.Attr) bool {
		// Проверяем ключ (приведи к нижнему регистру для надежности)
		key := strings.ToLower(a.Key)

		if _, ok := sensitiveKeys[key]; ok {
			newRecord.AddAttrs(slog.String(a.Key, "[REDACTED]"))
		} else {
			// Если ключ хороший, добавляем как есть
			newRecord.AddAttrs(a)
		}

		return true // Продолжаем итерацию
	})

	// 3. Передаем НОВУЮ (чистую) запись следующему хендлеру
	return s.Next.Handle(ctx, newRecord)
}

// Enabled reports whether the log level is allowed by the inner handler, preventing unnecessary masking operations.
func (s SecurityHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return s.Next.Enabled(ctx, l)
}

// WithAttrs returns a new SecurityHandler with redacted attributes to ensure safety in sub-loggers.
func (s SecurityHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cleanAttr := make([]slog.Attr, len(attrs))

	for i, a := range attrs {
		key := strings.ToLower(a.Key)
		if _, ok := sensitiveKeys[key]; ok {
			cleanAttr[i] = slog.String(a.Key, "[REDACTED]")

			a.Value = slog.StringValue("[REDACTED]")
		} else {
			// Если ключ хороший, добавляем как есть
			cleanAttr[i] = a
		}
	}

	return s.Next.WithAttrs(cleanAttr)
}

// WithGroup wraps the inner handler's group to maintain security masking in sub-loggers.
func (s SecurityHandler) WithGroup(name string) slog.Handler {
	return s.Next.WithGroup(name)
}
