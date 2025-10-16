-- for dashboard/dashboard-ui mock exchanges endpoint to create columns in dashboard-ui/src/table-pairs.tsx
CREATE TABLE exchanges (
    id TEXT PRIMARY KEY,
    name TEXT
);

-- for dashboard/dashboard-ui mock pairs endpoint to create rows in dashboard-ui/src/table-pairs.tsx
-- for fetcher fetched pairs 
CREATE TABLE pairs (
    id SERIAL PRIMARY KEY,
    name TEXT
);

insert into pairs (name) values ('btcusdt'),('bchusdt'),('ethusdt'),('adausdt'),('xrpusdt'),('solusdt'),('dogeusdt');