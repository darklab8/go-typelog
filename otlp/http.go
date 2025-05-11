package otlp

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func NewHTTPHandler(
	mux *http.ServeMux,
	register func(handleFunc func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request))),
) http.Handler {
	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	register(handleFunc)

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(spanNameFromPattern(mux), "/")
	return handler
}

// Opentelemetry not able to do it natively because they support go 1.22, and this thing needs 1.23
// https://github.com/open-telemetry/opentelemetry-go-contrib/issues/6193
// spanNameFromPattern is a simple middleware that sets the name of the span in the request context to the pattern used
// to match this request.
func spanNameFromPattern(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call handler first, so http.ServeMux can populate r.Pattern
		next.ServeHTTP(w, r)
		// Set span name after the fact. As long as this middleware is used within otelhttp.Handler, the span should
		// still be open and thus renameable.

		// Log.Info("serving from http", typelog.Int("result", roll))

		trace.SpanFromContext(r.Context()).SetName(r.Pattern)
	})
}
