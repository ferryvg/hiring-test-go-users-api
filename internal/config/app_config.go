package config

import (
	"sync"
	"time"
)

type AppConfig struct {
	Mysql MysqlConfig `mapstructure:"mysql"`
	Jwt   JwtConfig   `mapstructure:"jwt"`
}

type MysqlConfig struct {
	Nodes    *ClusterNodeList
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type JwtConfig struct {
	TTL time.Duration `mapstructure:"ttl"`
}

type ClusterNode struct {
	Host string
	Port int
}

func NewClusterNode(host string, port int) *ClusterNode {
	return &ClusterNode{
		Host: host,
		Port: port,
	}
}

type ClusterNodeList struct {
	nodes []*ClusterNode
	mu    *sync.RWMutex
}

func NewClusterNodeList() *ClusterNodeList {
	return &ClusterNodeList{
		mu: new(sync.RWMutex),
	}
}

func (l *ClusterNodeList) Get() []*ClusterNode {
	if l == nil {
		return nil
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	return append(l.nodes[0:0], l.nodes...)
}

func (l *ClusterNodeList) Set(nodes []*ClusterNode) {
	if l == nil {
		return
	}

	l.mu.Lock()

	l.nodes = append(l.nodes[0:0], nodes...)

	l.mu.Unlock()
}
