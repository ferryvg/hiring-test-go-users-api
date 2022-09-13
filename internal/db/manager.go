package db

import (
	"github.com/jmoiron/sqlx"
)

// Manager implements database connection manager.
type Manager interface {
	// Init initialize connection manager
	Init() error

	// Shutdown connection manager and clean up used resources.
	Shutdown()

	// GetDB returns connection for database.
	GetDB() (*sqlx.DB, error)
}

type managerImpl struct {
	connList ConnList
	resolver Resolver
	balancer Balancer
}

// NewManager creates database connection manager.
func NewManager(connList ConnList, resolver Resolver, balancer Balancer) Manager {
	return &managerImpl{
		connList: connList,
		resolver: resolver,
		balancer: balancer,
	}
}

func (m *managerImpl) Init() error {
	return m.resolver.Init()
}

func (m *managerImpl) Shutdown() {
	m.resolver.Shutdown()
	m.connList.Close()
}

func (m *managerImpl) GetDB() (*sqlx.DB, error) {
	return m.balancer.Next(m.connList.Connections())
}
