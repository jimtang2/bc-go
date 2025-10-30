package db

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("db.url", "postgres://postgres:password@localhost/bc?sslmode=disable")
}

var conn *sql.DB

func db() *sql.DB {
	if conn == nil {
		var err error
		conn, err = sql.Open("postgres", viper.GetString("db.url"))
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	return conn
}

func TrackedPairsAll() (pairs map[string][]string, err error) {
	var b []byte
	err = db().QueryRow(`select jsonb_object_agg(exchange, names) from (select exchange, array_to_json(array_agg(name)) as names from pairs group by exchange) sub`).Scan(&b)
	if err != nil {
		log.Println(err)
		return pairs, err
	}
	err = json.Unmarshal(b, &pairs)
	return pairs, err
}

func TrackedPairs(exchange string) (pairs []string) {
	var b []byte
	if err := db().QueryRow(`select array_to_json(array_agg(name)) from pairs where exchange=$1`, exchange).Scan(&b); err != nil {
		log.Println(err)
		return pairs
	}
	if err := json.Unmarshal(b, &pairs); err != nil {
		log.Println(err)
		return pairs
	}
	return pairs
}

func Exchanges() []Exchange {
	v := []Exchange{}
	rows, err := db().Query(`select id, name, taker_fee, maker_fee from exchanges`)
	if err != nil {
		return v
	}
	defer rows.Close()
	for rows.Next() {
		e := Exchange{}
		if err := rows.Scan(&e.ID, &e.Name, &e.TakerFee, &e.MakerFee); err != nil {
			return v
		}
		v = append(v, e)
	}
	return v
}

type Exchange struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	TakerFee float64 `json:"taker_fee"`
	MakerFee float64 `json:"maker_fee"`
}
