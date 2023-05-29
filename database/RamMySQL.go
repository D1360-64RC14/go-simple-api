package database

import (
	"database/sql"

	_ "github.com/proullon/ramsql/driver" // Needed to ramsql work

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/interfaces"
)

// RamMySQL implements Database
var _ interfaces.Database = (*RamMySQL)(nil)

type RamMySQL struct {
	settings *config.Database
	database *sql.DB
}

func NewRamMySQL(databaseSettings *config.Database) (interfaces.Database, error) {
	ramMySQL := &RamMySQL{
		settings: databaseSettings,
	}

	err := ramMySQL.setup()
	if err != nil {
		return nil, err
	}

	return ramMySQL, nil
}

func (d *RamMySQL) setup() (err error) {
	d.database, err = sql.Open("ramsql", d.Settings().DBName)
	if err != nil {
		return err
	}

	return nil
}

func (d RamMySQL) Settings() *config.Database {
	return d.settings
}

func (d RamMySQL) DB() *sql.DB {
	return d.database
}
