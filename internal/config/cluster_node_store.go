package config

import "sync"

type ClusterNodeStore interface {
	Get() []*ClusterNode
	Set(nodes []*ClusterNode)
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

type clusterNodeStoreImpl struct {
	nodes []*ClusterNode
	mu    *sync.RWMutex
}

func NewClusterNodeStore() ClusterNodeStore {
	return &clusterNodeStoreImpl{
		mu: new(sync.RWMutex),
	}
}

func (l *clusterNodeStoreImpl) Get() []*ClusterNode {
	if l == nil {
		return nil
	}

	l.mu.RLock()
	defer l.mu.RUnlock()

	return append(l.nodes[0:0], l.nodes...)
}

func (l *clusterNodeStoreImpl) Set(nodes []*ClusterNode) {
	if l == nil {
		return
	}

	l.mu.Lock()

	l.nodes = append(l.nodes[0:0], nodes...)

	l.mu.Unlock()
}
