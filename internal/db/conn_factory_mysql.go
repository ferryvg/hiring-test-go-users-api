package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" //nolint:gci
	"github.com/jmoiron/sqlx"
)

type mysqlConnFactory struct {
	database string
	username string
	password string
}

// NewMysqlConnFactory creates connections builder for MySQL.
func NewMysqlConnFactory(database, username, password string) ConnFactory {
	return &mysqlConnFactory{
		database: database,
		username: username,
		password: password,
	}
}

func (f *mysqlConnFactory) Create(node string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&loc=Local",
		f.username,
		f.password,
		node,
		f.database,
	)

	return sqlx.Connect("mysql", connStr)
}
