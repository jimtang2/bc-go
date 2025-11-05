# Arbitrage

## Description

```
git clone https://github.com/jimtang2/bc-go.git 
cd bc-go && docker compose up 
```

Visit `http://localhost:9000`

## System Components

1. Kafka
2. Postgres
3. AKHQ dashboard (optional)
4. bc_stream
5. bc_dashboard

## Basic Structure

| Process             |                                               |
| ------------------- | --------------------------------------------- |
| 1. Kafka            |                                               |
| 2. Postgres         |                                               |
| 3. AKHQ             |                                               |
| 4. bc_stream        | `cmd/stream`                                  |
| 5. bc_dashboard     | `cmd/dashboard`, `client`                     |

