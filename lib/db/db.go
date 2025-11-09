package db

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/jimtang2/bc-go/pkg/pb/v1"
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

func GetTrackedPairs() (pairs map[string][]string, err error) {
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

func Exchanges() []*pb.Exchange {
	v := []*pb.Exchange{}
	rows, err := db().Query(`select id, name, taker_fee, maker_fee from exchanges`)
	if err != nil {
		return v
	}
	defer rows.Close()
	for rows.Next() {
		e := pb.Exchange{}
		if err := rows.Scan(&e.Id, &e.Name, &e.TakerFee, &e.MakerFee); err != nil {
			return v
		}
		v = append(v, &e)
	}
	return v
}

func ExchangeFees() map[string]float64 {
	b := []byte{}
	err := db().QueryRow(`select jsonb_object_agg(id, taker_fee) from exchanges`).Scan(&b)
	v := map[string]float64{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		log.Println(err)
	}
	return v
}

// func InsertAlpha(a Alpha) error {
// 	_, err := db().Exec(`insert into alpha (id, pair, spread, volume, ask_exchange, ask_price, ask_size, ask_time, ask_fee, bid_exchange, bid_price, bid_size, bid_time, bid_fee, profit, timestamp) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) on conflict do nothing`, a.ID, a.Pair, a.Spread, a.Volume, a.AskExchange, a.AskPrice, a.AskSize, a.AskTime, a.AskFee, a.BidExchange, a.BidPrice, a.BidSize, a.BidTime, a.BidFee, a.Profit, a.Timestamp)
// 	return err
// }

// func AlphaItems() ([]byte, error) {
// 	b := []byte{}
// 	err := db().QueryRow(`select json_build_object('items', coalesce(json_agg(row_to_json(t)), '[]'), 'period', (select json_build_object('count', count(*), 'start', min(timestamp), 'end', max(timestamp), 'total', coalesce(sum(profit), 0.0)) from alpha)
//     ) as response from (select id as id, ask_exchange as ax, ask_price as ap, ask_size as as, ask_fee as af, ask_time as at, bid_exchange as bx, bid_price as bp, bid_size as bs, bid_fee as bf, bid_time as bt, pair as p, spread as s, volume as v, profit as pl, timestamp as ts from alpha where profit > 0.01 order by timestamp desc) t;`).Scan(&b)
// 	return b, err
// }

func InsertMatch(m *pb.Match) error {
	_, err := db().Exec(`
		INSERT INTO matches (
			id, pair,
			ask_exchange, ask_price, ask_size, ask_time, ask_fee_rate,
			bid_exchange, bid_price, bid_size, bid_time, bid_fee_rate,
			timestamp,
			calc_volume, calc_price_diff, calc_price_avg,
			calc_spread, calc_spread_pct, calc_bid_fee, calc_ask_fee, calc_profit_loss
		) VALUES (
			$1, $2,
			$3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13,
			$14, $15, $16,
			$17, $18, $19, $20, $21
		) ON CONFLICT (id) DO NOTHING`,
		// Match fields
		m.GetId(), m.GetPair(),
		m.GetAskExchange(), m.GetAskPrice(), m.GetAskSize(), m.GetAskTime(), m.GetAskFeeRate(),
		m.GetBidExchange(), m.GetBidPrice(), m.GetBidSize(), m.GetBidTime(), m.GetBidFeeRate(),
		m.GetTimestamp(),
		// Calculations
		m.GetCalculations().GetVolume(),
		m.GetCalculations().GetPriceDiff(),
		m.GetCalculations().GetPriceAvg(),
		m.GetCalculations().GetSpread(),
		m.GetCalculations().GetSpreadPct(),
		m.GetCalculations().GetBidFee(),
		m.GetCalculations().GetAskFee(),
		m.GetCalculations().GetProfitLoss(),
	)
	return err
}

func ProfitableMatches(limit, offset int) (*pb.ProfitableMatchesResponse, error) {
	if limit == 0 {
		limit = 1000
	}

	var resp pb.ProfitableMatchesResponse

	query := `
	WITH filtered AS (
	    SELECT 
	        id,
	        pair,
	        ask_exchange, ask_price, ask_size, ask_time, ask_fee_rate,
	        bid_exchange, bid_price, bid_size, bid_time, bid_fee_rate,
	        timestamp,
	        calc_volume, calc_price_diff, calc_price_avg,
	        calc_spread, calc_spread_pct, calc_bid_fee, calc_ask_fee, calc_profit_loss
	    FROM matches
	    WHERE calc_profit_loss > 0.1
	),
	paged AS (
	    SELECT row_to_json(t)::TEXT AS match_json
	    FROM (
	        SELECT
	            id,
	            pair,
	            ask_exchange,
	            ask_price,
	            ask_size,
	            ask_time,
	            ask_fee_rate,
	            bid_exchange,
	            bid_price,
	            bid_size,
	            bid_time,
	            bid_fee_rate,
	            timestamp,
	            (
	                SELECT row_to_json(c)
	                FROM (
	                    SELECT
	                        calc_volume AS volume,
	                        calc_price_diff AS price_diff,
	                        calc_price_avg AS price_avg,
	                        calc_spread AS spread,
	                        calc_spread_pct AS spread_pct,
	                        calc_bid_fee AS bid_fee,
	                        calc_ask_fee AS ask_fee,
	                        calc_profit_loss AS profit_loss
	                ) c
	            ) AS calculations
	        FROM filtered
	        ORDER BY timestamp DESC
	        LIMIT $1 OFFSET $2
	    ) t
	),
	agg AS (
	    SELECT
	        COUNT(*) AS total_count,
	        MIN(timestamp) AS min_ts,
	        COALESCE(SUM(calc_profit_loss), 0.0) AS total_profit
	    FROM filtered
	)
	SELECT JSON_BUILD_OBJECT(
	    'matches', COALESCE((SELECT JSON_AGG(match_json::json) FROM paged), '[]'),
	    'limit', $1,
	    'offset', $2,
	    'count', (SELECT total_count FROM agg),
	    'period_start', (SELECT min_ts FROM agg),
	    'period_profit', (SELECT total_profit FROM agg)
	)::TEXT;
	`

	var jsonResp []byte
	err := db().QueryRow(query, limit, offset).Scan(&jsonResp)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonResp, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
