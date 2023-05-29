package database

import (
	"database/sql"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/interfaces"
)

// MySQL implements Database
var _ interfaces.Database = (*MySQL)(nil)

type MySQL struct {
	settings *config.Database
	database *sql.DB
}

func NewMySQL(databaseSettings *config.Database) (interfaces.Database, error) {
	mySql := &MySQL{
		settings: databaseSettings,
	}

	err := mySql.setup()
	if err != nil {
		return nil, err
	}

	return mySql, nil
}

func (d *MySQL) setup() (err error) {
	d.database, err = sql.Open("mysql", d.Settings().DBName) // TODO: connect to MySQL database
	if err != nil {
		return err
	}

	return nil
}

func (d MySQL) Settings() *config.Database {
	return d.settings
}

func (d MySQL) DB() *sql.DB {
	return d.database
}
