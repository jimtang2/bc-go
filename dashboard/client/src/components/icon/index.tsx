import BinanceIcon from "./binance"
import BitfinexIcon from "./bitfinex"
import CoinbaseIcon from "./coinbase"
import KrakenIcon from "./kraken"
import OKXIcon from "./okx"

export default function ExchangeIcon({ exchange }: { exchange: string; }) {
	switch (exchange) {
	case "binance":
		return <BinanceIcon />
	case "bitfinex":
		return <BitfinexIcon />
	case "coinbase":
		return <CoinbaseIcon />
	case "kraken":
		return <KrakenIcon />
	case "okx":
		return <OKXIcon />
	default:
		return <></>
	}
}