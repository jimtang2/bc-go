
## Get Started

```bash
git clone https://github.com/jimtang2/bc-go.git
cd bc-go && docker compose up
```

This starts: 
- `bc-stream`: multiple streams of crypto pairs tickers from several CEX's (at this time, Binance, OKX, Kraken, Coinbase, Bitfinex)
- `bc-alpha`: stream processor to compute spread across exchanges
- `bc-dashboard`: visualization dashboard  

