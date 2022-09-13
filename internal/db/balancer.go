package db

import (
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
)

var ErrNoNodes = errors.New("no nodes")

// Balancer implement DB connections round-robin balancer.
type Balancer interface {

	// Next returns connection that should be used on next call.
	Next(conns []*sqlx.DB) (*sqlx.DB, error)
}

type roundRobinBalancer struct {
	next int
	mu   *sync.Mutex
}

// NewRoundRobinBalancer creates load balancer that use round-robin algorithm to spread work between nodes.
func NewRoundRobinBalancer() Balancer {
	return &roundRobinBalancer{mu: new(sync.Mutex)}
}

func (b *roundRobinBalancer) Next(conns []*sqlx.DB) (*sqlx.DB, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	connNum := len(conns)
	if connNum == 0 {
		return nil, ErrNoNodes
	}

	b.next = (b.next + 1) % connNum

	return conns[b.next], nil
}
