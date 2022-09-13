package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/consul"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/http"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/logger"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/metrics"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/saas"
	"github.com/ferryvg/hiring-test-go-users-api/internal/core/sd"
	"github.com/ferryvg/hiring-test-go-users-api/internal/dal"
	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
	server "github.com/ferryvg/hiring-test-go-users-api/internal/transport/http"
	"github.com/sirupsen/logrus"
)

func main() {
	app := core.NewApp()

	app.Set("svc.config_file", flag.String("config", "", "Configuration file to use"))

	app.Register(new(saas.Provider))
	app.Register(new(logger.Provider))
	app.Register(new(consul.Provider))
	app.Register(new(sd.Provider))
	app.Register(http.NewProvider(":6999"))
	app.Register(new(metrics.Provider))

	app.Register(new(config.Provider))
	app.Register(new(db.Provider))
	app.Register(new(dal.Provider))
	app.Register(new(server.Provider))

	log := app.MustGet("logger").(logrus.FieldLogger)
	log.Info("Starting service...")

	err := app.Boot()
	if err != nil {
		log.Fatalf("Failed to start service: %s", err)
	}

	log.Info("Service started!")

	go func(app *core.App, log logrus.FieldLogger) {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)

		for {
			<-sig
			err := app.Reconfigure()
			if err != nil {
				log.Errorf("Failed to reconfigure app: %s", err)
			}
		}
	}(app, log)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Info("Shutting down the service...")
	app.Shutdown()
	log.Info("Service stopped!")
}
