package config

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	sspConfigFilename = "users-api"
	sspConfigFileType = "yaml"
)

//nolint:gochecknoglobals
var possibleConfigPaths = []string{
	"$HOME",
	"$HOME/.config/ssp",
	".",
	"./configs",
}

type Builder interface {
	Build(config string) (*AppConfig, error)
}

type builderImpl struct {
	viper  *viper.Viper
	logger logrus.FieldLogger
}

func NewBuilder(logger logrus.FieldLogger) Builder {
	return &builderImpl{viper: viper.New(), logger: logger}
}

func (b *builderImpl) Build(config string) (*AppConfig, error) {
	b.configureViper()

	err := b.applyConfigFile(config)
	if err != nil {
		return nil, err
	}

	var appConfig AppConfig

	err = b.viper.Unmarshal(&appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}

func (b *builderImpl) configureViper() {
	b.viper.SetEnvPrefix("APP")
	b.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	b.viper.AutomaticEnv()

	b.viper.SetDefault("mysql", map[string]interface{}{
		"database": "users_api",
		"username": "root",
		"password": "scout",
	})

	b.viper.SetDefault("jwt", map[string]interface{}{
		"ttl": "720h", // 30 days
	})
}

func (b *builderImpl) applyConfigFile(config string) error {
	var ignoreNotFound bool

	if config != "" {
		config, _ = filepath.Abs(config)

		b.viper.SetConfigFile(config)
	} else {
		ignoreNotFound = true

		for _, path := range possibleConfigPaths {
			b.viper.AddConfigPath(path)
		}

		b.viper.SetConfigName(sspConfigFilename)
		b.viper.SetConfigType(sspConfigFileType)
	}

	err := b.viper.ReadInConfig()
	if err == nil {
		b.logger.Infof("Loading configuration from %s", b.viper.ConfigFileUsed())

		return nil
	}

	var viErr viper.ConfigFileNotFoundError
	if ok := errors.As(err, &viErr); ok && ignoreNotFound {
		b.logger.WithError(err).Warn("Configuration file is not found, using defaults or env")

		return nil
	}

	return err
}
