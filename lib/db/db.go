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

func InsertAlpha(a Alpha) error {
	_, err := db().Exec(`insert into alpha (id, pair, spread, volume, ask_exchange, ask_price, ask_size, ask_time, ask_fee, bid_exchange, bid_price, bid_size, bid_time, bid_fee, profit, timestamp) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) on conflict do nothing`, a.ID, a.Pair, a.Spread, a.Volume, a.AskExchange, a.AskPrice, a.AskSize, a.AskTime, a.AskFee, a.BidExchange, a.BidPrice, a.BidSize, a.BidTime, a.BidFee, a.Profit, a.Timestamp)
	return err
}

type Alpha struct {
	ID          int64   `json:"i"` // offset id of 'alpha' topic
	Pair        string  `json:"p"`
	Spread      float64 `json:"s"`
	Volume      float64 `json:"v"`
	AskExchange string  `json:"ax"`
	AskPrice    float64 `json:"ap"`
	AskSize     float64 `json:"as"`
	AskTime     int64   `json:"at"`
	AskFee      float64 `json:"af"`
	BidExchange string  `json:"bx"`
	BidPrice    float64 `json:"bp"`
	BidSize     float64 `json:"bs"`
	BidTime     int64   `json:"bt"`
	BidFee      float64 `json:"bf"`
	Profit      float64 `json:"pl"`
	Timestamp   int64   `json:"ts"`
}

type AlphaResponse struct {
	Items  []Alpha `json:"items"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
	Period struct {
		Count int     `json:"count"`
		Start int64   `json:"start"`
		Total float64 `json:"total"`
	} `json:"period"`
}

func AlphaItems() ([]byte, error) {
	b := []byte{}
	err := db().QueryRow(`select json_build_object('items', coalesce(json_agg(row_to_json(t)), '[]'), 'period', (select json_build_object('count', count(*), 'start', min(timestamp), 'end', max(timestamp), 'total', coalesce(sum(profit), 0.0)) from alpha)
    ) as response from (select id as id, ask_exchange as ax, ask_price as ap, ask_size as as, ask_fee as af, ask_time as at, bid_exchange as bx, bid_price as bp, bid_size as bs, bid_fee as bf, bid_time as bt, pair as p, spread as s, volume as v, profit as pl, timestamp as ts from alpha order by timestamp desc) t;`).Scan(&b)
	return b, err
}
