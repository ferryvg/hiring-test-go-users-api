package dal_test

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

type DbManager struct {
	host     string
	port     int
	db       string
	user     string
	password string
}

// Creates database connection DbManager
func NewDbManager() *DbManager {
	return &DbManager{
		host:     "db",
		port:     3306,
		db:       "users_api",
		user:     "root",
		password: "scout",
	}
}

func (m *DbManager) Init() error {
	if db, ok := os.LookupEnv("APP_MYSQL_DATABASE"); ok {
		m.db = db
	}

	if user, ok := os.LookupEnv("APP_MYSQL_USERNAME"); ok {
		m.user = user
	}

	if password, ok := os.LookupEnv("APP_MYSQL_PASSWORD"); ok {
		m.password = password
	}

	return nil
}

func (m *DbManager) Shutdown() {
}

func (m *DbManager) GetDB() (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		m.user,
		m.password,
		m.host,
		m.port,
		m.db,
	)

	return sqlx.Open("mysql", connStr)
}
