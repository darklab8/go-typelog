package otlp

import (
	"context"
	"errors"
	stdlog "log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// SetupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTracerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return
}

func GetResource() *resource.Resource {
	res, err := resource.New(
		context.Background(),
		resource.WithFromEnv(),   // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
		resource.WithContainer(), // Discover and provide container information.
		resource.WithAttributes(attribute.String("environment", "dev")),
	)
	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		stdlog.Println(err) // Log non-fatal issues.
	} else if err != nil {
		stdlog.Fatalln(err) // The error may be fatal.
	}
	return res
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider() (*trace.TracerProvider, error) {
	traceExporter, err2 := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure())
	if err2 != nil {
		return nil, err2
	}
	_ = traceExporter

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(GetResource()),
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
	)
	return tracerProvider, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err2 := otlpmetrichttp.New(context.Background(), otlpmetrichttp.WithInsecure())
	if err2 != nil {
		return nil, err2
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(GetResource()),

		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}

func newLoggerProvider() (*log.LoggerProvider, error) {
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}
	_ = logExporter
	httplogExporter, err2 := otlploghttp.New(context.Background(), otlploghttp.WithInsecure())
	if err2 != nil {
		return nil, err2
	}
	_ = httplogExporter

	loggerProvider := log.NewLoggerProvider(
		log.WithResource(GetResource()),
		log.WithProcessor(log.NewBatchProcessor(httplogExporter)),
		// log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
