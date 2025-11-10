-- for dashboard/dashboard-ui mock exchanges endpoint to create columns in dashboard-ui/src/table-pairs.tsx
CREATE TABLE exchanges (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    taker_fee DECIMAL NOT NULL,
    maker_fee DECIMAL NOT NULL
);

INSERT INTO exchanges (id, name, taker_fee, maker_fee) VALUES
    ('binance',   'Binance',   0.10, 0.10),  -- 0.10%
    ('bitfinex',  'Bitfinex',  0.20, 0.10),  -- 0.20% taker, 0.10% maker
    ('uniswap',   'Uniswap',   0.30, 0.30),  -- ~0.30% (varies by pool)
    ('coinbase',  'Coinbase',  0.60, 0.40),  -- 0.60% taker, 0.40% maker
    ('kraken',    'Kraken',    0.26, 0.16),  -- 0.26% taker, 0.16% maker
    ('bybit',     'Bybit',     0.10, 0.06),  -- 0.10% taker, 0.06% maker
    ('okx',       'OKX',       0.10, 0.08),  -- 0.10% taker, 0.08% maker
    ('gemini',    'Gemini',    0.35, 0.25)   -- 0.35% taker, 0.25% maker
ON CONFLICT DO NOTHING;

-- for dashboard/dashboard-ui mock pairs endpoint to create rows in dashboard-ui/src/table-pairs.tsx
-- for fetcher fetched pairs 
CREATE TABLE pairs (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    exchange TEXT NOT NULL,
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

CREATE TABLE IF NOT EXISTS alpha (
    id              BIGINT PRIMARY KEY,
    pair            TEXT NOT NULL,
    spread          DOUBLE PRECISION NOT NULL,
    volume          DOUBLE PRECISION NOT NULL,
    ask_exchange    TEXT NOT NULL,
    ask_price       DOUBLE PRECISION NOT NULL,
    ask_size        DOUBLE PRECISION NOT NULL,
    ask_time        BIGINT NOT NULL,
    ask_fee         DOUBLE PRECISION NOT NULL,
    bid_exchange    TEXT NOT NULL,
    bid_price       DOUBLE PRECISION NOT NULL,
    bid_size        DOUBLE PRECISION NOT NULL,
    bid_time        BIGINT NOT NULL,
    bid_fee         DOUBLE PRECISION NOT NULL,
    profit          DOUBLE PRECISION NOT NULL,
    timestamp       BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS matches (
    id                  BIGINT PRIMARY KEY,
    pair                TEXT NOT NULL,             
    ask_exchange        TEXT NOT NULL,             
    ask_price           DOUBLE PRECISION NOT NULL, 
    ask_size            DOUBLE PRECISION NOT NULL, 
    ask_time            BIGINT NOT NULL,           
    ask_fee_rate        DOUBLE PRECISION NOT NULL, 
    bid_exchange        TEXT NOT NULL,             
    bid_price           DOUBLE PRECISION NOT NULL, 
    bid_size            DOUBLE PRECISION NOT NULL, 
    bid_time            BIGINT NOT NULL,           
    bid_fee_rate        DOUBLE PRECISION NOT NULL, 
    timestamp           BIGINT NOT NULL,           
    calc_volume         DOUBLE PRECISION NOT NULL, 
    calc_price_diff     DOUBLE PRECISION NOT NULL, 
    calc_price_avg      DOUBLE PRECISION NOT NULL, 
    calc_spread         DOUBLE PRECISION NOT NULL, 
    calc_spread_pct     DOUBLE PRECISION NOT NULL, 
    calc_bid_fee        DOUBLE PRECISION NOT NULL, 
    calc_ask_fee        DOUBLE PRECISION NOT NULL,  
    calc_profit_loss    DOUBLE PRECISION NOT NULL   
);