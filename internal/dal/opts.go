package dal

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Options struct {
	Ctx  context.Context
	Conn *sqlx.DB
	Tx   *sqlx.Tx
}

type ManagerOpt func(opts *Options)

func WithConn(conn *sqlx.DB) ManagerOpt {
	return func(opts *Options) {
		opts.Conn = conn
	}
}

func WithTx(tx *sqlx.Tx) ManagerOpt {
	return func(opts *Options) {
		opts.Tx = tx
	}
}
