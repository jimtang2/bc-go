# Arbitrage

## Get Started

```bash
git clone https://github.com/jimtang2/bc-go.git
cd bc-go && docker compose up
```

This starts: 
- `bc-stream`: multiple streams of crypto pairs tickers from several CEX's (at this time, Binance, OKX, Kraken, Coinbase, Bitfinex)
- `bc-alpha`: stream processor to compute spread across exchanges
- `bc-dashboard`: visualization dashboard  

## Development Notes

- tracked pairs are set in postgres (table `pairs`); refer to deploy/config/config.yml for connection details 
- new streams can be added by implementing `github.com/jimtang2/bc-go/lib/stream#Stream` interface and editing `cmd/stream/proxy.go`
- vite dev mode leaks memory 

## Next

- wallet connector