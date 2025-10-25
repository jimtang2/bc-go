package db

import (
	"encoding/json"
	"log"
)

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
