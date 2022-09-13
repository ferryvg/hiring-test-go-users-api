package http

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net"
	"os"
	"time"
)

const (
	// Setup address that server wil listen
	serverAddrEnvName = "HTTP_SERVER_ADDR"
	// Setup timeout for graceful shutdown
	shutdownTimeoutEnvName = "HTTP_SHUTDOWN_TIMEOUT"
)

const defaultShutdownTimeout = 30 * time.Second

type Provider struct {
	addr string
}

func NewProvider(addr string) *Provider {
	return &Provider{addr: addr}
}

func (p *Provider) Register(c core.Container) {
	p.registerAddr(c)
	p.registerShutdownTimeout(c)
	p.registerListener(c)
	p.registerRouters(c)
	p.registerServer(c)
}

func (p *Provider) Boot(c core.Container) error {
	listener := c.MustGet("http.listener").(net.Listener)
	server := c.MustGet("http.server").(*fasthttp.Server)

	go func() {
		_ = server.Serve(listener)
	}()

	return nil
}

func (p *Provider) Shutdown(c core.Container) {
	listener := c.MustGet("http.listener").(net.Listener)
	err := listener.Close()

	if err != nil {
		log := c.MustGet("logger").(logrus.FieldLogger)
		log.Warnf("Failed to gracefully shutdown HTTP server: %s", err)
	}
}

func (p *Provider) registerAddr(c core.Container) {
	c.Set("http.addr", func(c core.Container) interface{} {
		addr, exists := os.LookupEnv(serverAddrEnvName)
		if exists {
			return addr
		}

		return p.addr
	})
}

func (p *Provider) registerShutdownTimeout(c core.Container) {
	c.Set("http.shutdown_timeout", func(c core.Container) interface{} {
		timeoutStr, exists := os.LookupEnv(shutdownTimeoutEnvName)
		if exists {
			timeout, err := time.ParseDuration(timeoutStr)
			if err == nil {
				return timeout
			} else {
				log := c.MustGet("logger").(logrus.FieldLogger)
				log.Errorf(
					"Failed parse '%v' env variable. Value - %v. Error - %v\n",
					shutdownTimeoutEnvName,
					timeoutStr,
					err,
				)
			}
		}

		return defaultShutdownTimeout
	})
}

func (p *Provider) registerListener(c core.Container) {
	c.Set("http.listener", func(c core.Container) interface{} {
		addr := c.MustGet("http.addr").(string)
		timeout := c.MustGet("http.shutdown_timeout").(time.Duration)

		ln, err := NewListener(addr, timeout)
		if err != nil {
			logger := c.MustGet("logger").(logrus.FieldLogger)
			logger.Fatalf("Failed to start HTTP server: %s", err)
		}

		return ln
	})
}

func (p *Provider) registerRouters(c core.Container) {
	c.Set("http.router", func(c core.Container) interface{} {
		return fasthttprouter.New()
	})
}

func (p *Provider) registerServer(c core.Container) {
	c.Set("http.server", func(c core.Container) interface{} {
		router := c.MustGet("http.router").(*fasthttprouter.Router)
		logger := c.MustGet("logger").(logrus.FieldLogger)

		return &fasthttp.Server{
			Handler:        router.Handler,
			Logger:         logger,
			ReadBufferSize: 1 << 15,
		}
	})
}
