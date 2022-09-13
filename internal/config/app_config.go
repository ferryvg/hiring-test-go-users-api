package config

import (
	"time"
)

type AppConfig struct {
	Mysql MysqlConfig `mapstructure:"mysql"`
	Jwt   JwtConfig   `mapstructure:"jwt"`
}

type MysqlConfig struct {
	Nodes    ClusterNodeStore
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type JwtConfig struct {
	TTL time.Duration `mapstructure:"ttl"`
}
