import React from "react"
import BinanceIcon from "./binance"
import BitfinexIcon from "./bitfinex"
import CoinbaseIcon from "./coinbase"
import KrakenIcon from "./kraken"
import OKXIcon from "./okx"

const icons: Record<string, () => React.ReactElement> = {
	"binance": () => <BinanceIcon />,
	"bitfinex": () => <BitfinexIcon />,
	"coinbase": () => <CoinbaseIcon />,
	"kraken": () => <KrakenIcon />,
	"okx": () => <OKXIcon />,
}

export default icons