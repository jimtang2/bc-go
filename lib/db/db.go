package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

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
	rows, err := db().Query(`select id, name, taker_fee, maker_fee, chart_color from exchanges`)
	if err != nil {
		return v
	}
	defer rows.Close()
	for rows.Next() {
		e := pb.Exchange{}
		if err := rows.Scan(&e.Id, &e.Name, &e.TakerFee, &e.MakerFee, &e.ChartColor); err != nil {
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
	var resp pb.ProfitableMatchesResponse

	query := `
	with filtered as ( 
		select id, pair, ask_exchange, ask_price, ask_size, ask_time, ask_fee_rate, bid_exchange, bid_price, bid_size, bid_time, bid_fee_rate, timestamp, calc_volume, calc_price_diff, calc_price_avg, calc_spread, calc_spread_pct, calc_bid_fee, calc_ask_fee, calc_profit_loss 
		from matches 
		where calc_profit_loss > 0.1 
	),
	paged as ( 
		select row_to_json(t)::text as match_json from ( 
			select 
			id, pair, ask_exchange, ask_price, ask_size, ask_time, ask_fee_rate, bid_exchange, bid_price, bid_size, bid_time, bid_fee_rate, timestamp, 
				( select row_to_json(c) from ( 
						select calc_volume as volume, calc_price_diff as price_diff, calc_price_avg as price_avg, calc_spread as spread, calc_spread_pct as spread_pct, calc_bid_fee as bid_fee, calc_ask_fee as ask_fee, calc_profit_loss as profit_loss 
					) c ) as calculations
	    from filtered
	        order by timestamp desc
	        LIMIT $1 OFFSET $2
	    ) t
	),
	agg as (
	    select count(*) as total_count, min(timestamp) as min_ts, coalesce(sum(calc_profit_loss), 0.0) as total_profit from filtered
	)
	select json_build_object( 
		'matches', 
		coalesce((select json_agg(match_json::json) from paged), '[]'), 
		'limit', $1, 
		'offset', $2, 
		'count', (SELECT total_count from agg), 
		'period_start', (SELECT min_ts from agg), 
		'period_profit', (SELECT total_profit from agg)
	)::text;
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

func LogVisit(r *http.Request) error {
	ip := getClientIP(r)
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	queryString := r.URL.RawQuery
	query := `
		INSERT INTO visits (
			ip_address, user_agent, referer, path, method, protocol,
			host, query_string, remote_addr
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`
	_, err := db().Exec(query,
		ip,
		r.UserAgent(),
		r.Referer(),
		path,
		r.Method,
		r.Proto,
		r.Host,
		queryString,
		r.RemoteAddr,
	)
	return err
}

func getClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if real := r.Header.Get("X-Real-IP"); real != "" {
		return real
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if host != "" {
		return host
	}
	return r.RemoteAddr
}
