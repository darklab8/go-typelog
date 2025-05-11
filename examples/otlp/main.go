package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"math/rand"

	"github.com/darklab8/go-typelog/otlp"
	"github.com/darklab8/go-typelog/typelog"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

var (
	Log2 *typelog.Logger = typelog.NewLogger("infra", typelog.WithLogLevel(typelog.LEVEL_INFO))
	Log  *typelog.Logger = Log2.WithScope("infra.main")
)

const name = "github.com/darklab8/infra"

var (
	tracer  = otel.Tracer(name)
	meter   = otel.Meter(name)
	logger  = otelslog.NewLogger(name) // not catching stuff not going through it, stdout format is mess
	rollCnt metric.Int64Counter
)

func init() {
	var err error
	rollCnt, err = meter.Int64Counter("dice.rolls",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}

}

func rolldice(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End()

	Log.InfoCtx(ctx, "rolldice started")

	roll := 1 + rand.Intn(6)

	var msg string
	if player := r.PathValue("player"); player != "" {
		msg = fmt.Sprintf("%s is rolling the dice", player)
	} else {
		msg = "Anonymous player is rolling the dice"
	}
	if false {
		// logger with many minuses.
		logger.InfoContext(ctx, msg, "result", roll)
	}

	Log2.InfoCtx(ctx, msg, typelog.Int("result", roll))
	Log.ErrorCtx(ctx, msg, typelog.Bool("is_attr", true))

	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)
	rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

func run() (err error) {
	Log.Info("started run")
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	otelShutdown, err := otlp.SetupOTelSDK(ctx) // Set up OpenTelemetry.
	if err != nil {
		return
	}
	defer func() { // Handle shutdown properly so nothing leaks.
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler: otlp.NewHTTPHandler(http.NewServeMux(), func(handleFunc func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request))) {
			// Register handlers.
			handleFunc("/rolldice/", rolldice)
			handleFunc("/rolldice/{player}", rolldice)
		}),
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}
