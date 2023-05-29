package interfaces

import (
	"database/sql"

	"github.com/d1360-64rc14/simple-api/config"
)

type Database interface {
	Settings() *config.Database
	DB() *sql.DB
	Close() error
}
