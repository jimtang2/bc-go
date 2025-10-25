-- for dashboard/dashboard-ui mock exchanges endpoint to create columns in dashboard-ui/src/table-pairs.tsx
CREATE TABLE exchanges (
    id TEXT PRIMARY KEY,
    name TEXT
);

insert into exchanges (id, name) values
    ('binance','Binance'),
    ('bitfinex','Bitfinex'),
    ('uniswap','Uniswap'),
    ('coinbase','Coinbase'),
    ('kraken','Kraken'),
    ('bybit','Bybit'),
    ('okx','OKX'),
    ('gemini','Gemini')
on conflict do nothing;


-- for dashboard/dashboard-ui mock pairs endpoint to create rows in dashboard-ui/src/table-pairs.tsx
-- for fetcher fetched pairs 
CREATE TABLE pairs (
    id SERIAL PRIMARY KEY,
    name TEXT,
    exchange TEXT,
    UNIQUE(name, exchange)
);

insert into pairs (exchange, name) values 
    ('binance','adausdt'),
    ('binance','adausdc'),
    ('binance','bchusdt'),
    ('binance','bchusdc'),
    ('binance','bnbusdt'),
    ('binance','bnbusdc'),
    ('binance','btcusdt'),
    ('binance','btcusdc'),
    ('binance','dogeusdt'),
    ('binance','dogeusdc'),
    ('binance','ethusdt'),
    ('binance','ethusdc'),
    ('binance','solusdt'),
    ('binance','solusdc'),
    ('binance','trxsdt'),
    ('binance','trxusdc'),
    ('binance','xrpusdt'),
    ('binance','xrpusdc'),
    ('bitfinex','ADAUSD'),
    ('bitfinex','BTCUSD'),
    ('bitfinex','DOGE:USD'),
    ('bitfinex','ETHUSD'),
    ('bitfinex','SOLUSD'),
    ('bitfinex','TRXUSD'),
    ('bitfinex','XRPUSD'),
    ('coinbase','ADA-USDT'),
    ('coinbase','BTC-USDT'),
    ('coinbase','DOGE-USDT'),
    ('coinbase','ETH-USDT'),
    ('coinbase','SOL-USDT'),
    ('coinbase','XRP-USDT'),
    ('kraken', 'ADA/USDT'),
    ('kraken', 'ADA/USDC'),
    ('kraken', 'BCH/USDT'),
    ('kraken', 'BCH/USDC'),
    ('kraken', 'BNB/USDT'),
    ('kraken', 'BNB/USDC'),
    ('kraken', 'BTC/USDT'),
    ('kraken', 'BTC/USDC'),
    ('kraken', 'DOGE/USDT'),
    ('kraken', 'DOGE/USDC'),
    ('kraken', 'ETH/USDT'),
    ('kraken', 'ETH/USDC'),
    ('kraken', 'SOL/USDT'),
    ('kraken', 'SOL/USDC'),
    ('kraken', 'TRX/SDT'),
    ('kraken', 'TRX/USDC'),
    ('kraken', 'XRP/USDT'),
    ('kraken', 'XRP/USDC'),
    ('okx', 'ADA-USDT'),
    ('okx', 'ADA-USDC'),
    ('okx', 'BCH-USDT'),
    ('okx', 'BCH-USDC'),
    ('okx', 'BNB-USDT'),
    ('okx', 'BNB-USDC'),
    ('okx', 'BTC-USDT'),
    ('okx', 'BTC-USDC'),
    ('okx', 'DOGE-USDT'),
    ('okx', 'DOGE-USDC'),
    ('okx', 'ETH-USDT'),
    ('okx', 'ETH-USDC'),
    ('okx', 'SOL-USDT'),
    ('okx', 'SOL-USDC'),
    ('okx', 'TRX-USDT'),
    ('okx', 'TRX-USDC'),
    ('okx', 'XRP-USDT'),
    ('okx', 'XRP-USDC')
on conflict do nothing;