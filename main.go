package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	fmt.Println("default", strings.Repeat("-", 20))
	defaultLogger()

	fmt.Println("json", strings.Repeat("-", 20))
	jsonLogger()

	fmt.Println("context", strings.Repeat("-", 20))
	ctxLogger()
}

func defaultLogger() {
	logger := slog.Default()
	logThings(context.Background(), logger)
}

func ctxLogger() {
	logger := slog.New(&contextLogger{slog.NewJSONHandler(os.Stdout)})
	ctx := appendAttrs(context.Background(), slog.String("traceID", "someid"))
	logThings(ctx, logger)
}

func jsonLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout))
	logThings(context.Background(), logger)
}

func logThings(ctx context.Context, logger *slog.Logger) {
	// Ugly format
	logger.InfoCtx(ctx, "Hello, world", "GOOS", runtime.GOOS)

	// Structured Format
	logger.InfoCtx(ctx, "Hello, mars", slog.Time("time", time.Now()))

	// Include values
	logger = logger.With(slog.Int("anumber", 5))
	logger.InfoCtx(ctx, "Log with number")

	// Grouped Values
	logger.InfoCtx(ctx, "runtime info", slog.Group("runtime",
		slog.String("GOOS", runtime.GOOS),
		slog.String("GOARCH", runtime.GOARCH),
	))

	// Child logger
	childLogger := logger.WithGroup("child")
	childLogger.InfoCtx(ctx, "some log", "a", 1, "b", 2)
}
