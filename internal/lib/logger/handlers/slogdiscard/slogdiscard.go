package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil // Ignore log writing
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h // return the same, because don`t have attr for save
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h // return the same, because don`t have group for save
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false // always return false, because ignored log writing
}
