package database

import "database/sql"

var DB *sql.DB

func SetDB(d *sql.DB) {
	DB = d
}
