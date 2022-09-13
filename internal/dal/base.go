package dal

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/ferryvg/hiring-test-go-users-api/internal/db"
)

type baseManager struct {
	manager db.Manager
	logger  logrus.FieldLogger
}

func (m *baseManager) resolveOptions(ctx context.Context, optsList []ManagerOpt) (*Options, error) {
	opts := &Options{
		Ctx: ctx,
	}

	for _, opt := range optsList {
		opt(opts)
	}

	if err := m.checkConn(opts); err != nil {
		return nil, err
	}

	return opts, nil
}

func (m *baseManager) checkConn(opts *Options) error {
	if opts.Conn != nil {
		return nil
	}

	conn, err := m.manager.GetDB()
	if err != nil {
		return fmt.Errorf("failed to create db connection: %w", err)
	}

	opts.Conn = conn

	return nil
}
