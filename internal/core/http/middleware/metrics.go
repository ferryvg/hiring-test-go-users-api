package middleware

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
)

var (
	serverRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "handled_total",
			Help:      "Total number of handled HTTP requests.",
		},
		[]string{"http_endpoint", "http_method", "http_code"},
	)

	serverResponseHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "handling_seconds",
			Help:      "Histogram of HTTP response latency (seconds)",
			Buckets:   []float64{.050, .100, .200, .300, .400, .500, 1},
		},
		[]string{"http_endpoint", "http_method"},
	)
)

func init() {
	prometheus.MustRegister(serverRequestCounter)
	prometheus.MustRegister(serverResponseHistogram)
}

func Metrics(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()
		endpoint := string(ctx.Request.URI().Path())
		method := string(ctx.Request.Header.Method())

		next(ctx)

		code := strconv.Itoa(ctx.Response.StatusCode())

		serverRequestCounter.WithLabelValues(endpoint, method, code).Inc()
		serverResponseHistogram.WithLabelValues(endpoint, method).Observe(time.Since(startTime).Seconds())
	}
}
