package main

import (
    "io"
    "log/slog"
    "os"
    "strings"
)

var Logger *slog.Logger

func SetupLogger(level, format, file string, noColor bool) *slog.Logger {
    var w io.Writer = os.Stdout
    if strings.TrimSpace(file) != "" {
        f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err == nil {
            w = f
        } else {
            w = os.Stdout
        }
    }

    var lvl slog.Level
    switch strings.ToLower(strings.TrimSpace(level)) {
    case "debug":
        lvl = slog.LevelDebug
    case "warn":
        lvl = slog.LevelWarn
    case "error":
        lvl = slog.LevelError
    default:
        lvl = slog.LevelInfo
    }

    opts := &slog.HandlerOptions{Level: lvl}

    var handler slog.Handler
    switch strings.ToLower(strings.TrimSpace(format)) {
    case "json":
        handler = slog.NewJSONHandler(w, opts)
    default:
        handler = slog.NewTextHandler(w, opts)
    }

    return slog.New(handler)
}