package metrics

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// Service provider that register components required for metrics collection
//
// It's register:
//    - middleware that collect metrics for gRPC server calls
//    - HTTP route that report collected metrics
type Provider struct{}

// Register components in service container
func (p *Provider) Register(c core.Container) {
	p.registerHttpRoute(c)
}

func (p *Provider) registerHttpRoute(c core.Container) {
	c.MustExtend("http.router", func(old interface{}, c core.Container) interface{} {
		routes := old.(*fasthttprouter.Router)
		routes.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

		return routes
	})
}
