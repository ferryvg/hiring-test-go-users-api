package dal_test

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
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
		host:     "localhost",
		port:     3306,
		db:       "test_db",
		user:     "root",
		password: "scout",
	}
}

func (m *DbManager) Init() error {
	if host, ok := os.LookupEnv("MYSQL_HOST"); ok {
		m.host = host
	}

	if portStr, ok := os.LookupEnv("MYSQL_PORT"); ok {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			m.port = port
		}
	}

	if db, ok := os.LookupEnv("MYSQL_DB"); ok {
		m.db = db
	}

	if user, ok := os.LookupEnv("MYSQL_USER"); ok {
		m.user = user
	}

	if password, ok := os.LookupEnv("MYSQL_PASSWORD"); ok {
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
