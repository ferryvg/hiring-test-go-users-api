package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core"
	"github.com/sirupsen/logrus"
)

const (
	levelEnvName    = "LOG_LEVEL"
	defaultLogLevel = logrus.InfoLevel

	formatEnvName = "LOG_FORMAT"
	defaultFormat = "text"
)

type Provider struct{}

func (p *Provider) Reconfigure(container core.Container) {
	log := container.MustGet("logger").(*logrus.Logger)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Errorf("Failed to detect binary path: %s", err)
		return
	}

	path := dir + string(os.PathSeparator) + levelEnvName
	bytes, err := ioutil.ReadFile(path)
	if err == nil {
		lvl, err := logrus.ParseLevel(string(bytes))
		if err != nil {
			log.Errorf("Failed to parse content of %s. Value - %v. Error - %v", levelEnvName, string(bytes), err)
		} else {
			log.SetLevel(lvl)
		}
	} else {
		log.Warn(err)
	}

	path = dir + string(os.PathSeparator) + formatEnvName
	bytes, err = ioutil.ReadFile(path)
	if err == nil {
		formatter, err := p.parseLogFormat(string(bytes))
		if err != nil {
			log.Errorf("Failed to parse content of %s. Value - %v. Error - %v", formatEnvName, string(bytes), err)
		} else {
			log.Formatter = formatter
		}
	} else {
		log.Warn(err)
	}
}

func (p *Provider) Register(container core.Container) {
	p.registerLogLevel(container)
	p.registerLogFormat(container)

	container.Set("logger", func(c core.Container) interface{} {
		level := c.MustGet("logger.level").(logrus.Level)
		formatter := c.MustGet("logger.format").(logrus.Formatter)

		logger := logrus.New()
		logger.Formatter = formatter
		logger.SetLevel(level)
		return logger
	})
}

func (p *Provider) registerLogLevel(c core.Container) {
	c.Set("logger.level", func(c core.Container) interface{} {
		env, exist := os.LookupEnv(levelEnvName)

		if exist && env != "" {
			lvl, err := logrus.ParseLevel(env)
			if err != nil {
				// we can't use c.MustGet("logger") here
				logrus.Errorf(
					"Failed parse '%v' env variable. Value - %v. Error - %v",
					levelEnvName,
					env,
					err,
				)
			} else {
				return lvl
			}
		}

		return defaultLogLevel
	})
}

func (p *Provider) registerLogFormat(c core.Container) {
	c.Set("logger.format", func(c core.Container) interface{} {
		format, exist := os.LookupEnv(formatEnvName)
		if exist && format != "" {
			formatter, err := p.parseLogFormat(format)
			if err != nil {
				logrus.Errorf("Failed parse '%v' env variable: %s. Use 'text' as default", formatEnvName, err)
			} else {
				return formatter
			}
		}

		defaultFormatter, _ := p.parseLogFormat(defaultFormat)

		return defaultFormatter
	})
}

func (p *Provider) parseLogFormat(str string) (logrus.Formatter, error) {
	format := strings.ToLower(str)

	switch format {
	case "json":
		return &logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano}, nil
	case "text":
		return &logrus.TextFormatter{TimestampFormat: time.RFC3339Nano}, nil
	default:
		return nil, fmt.Errorf("incorrect format: '%v'. Allowed options: 'text','json'", format)
	}
}
