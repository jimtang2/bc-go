// webhook$$bc-go;dashboard/server/db.go;grok$$
package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db = &DB{}

type DB struct {
	db *sql.DB
}

func (d *DB) lists() ([]byte, error) {
	var (
		b   []byte
		err error
	)
	if d.db == nil {
		d.db, err = sql.Open("postgres", viper.GetString("db.url"))
		if err != nil {
			return b, err
		}
	}
	err = d.db.QueryRow(`select jsonb_build_object(
  'exchanges', (select array_to_json(array_agg(id)) from exchanges),
  'pairs', (select array_to_json(array_agg(name)) from pairs))`).Scan(&b)
	return b, err
}
