# Arbitrage

## Description

Application code for the arbitrage system

### Dependencies

1. `github.com/IBM/sarama`
2. `github.com/lib/pq`

## Project Structure

| Process       | Topic Output | Packages              | Dep Services |
| ------------- | ------------ | --------------------- | ------------ |
| `cmd/tickers` | `tickers`    | `lib/stream` `lib/db` | `kafka` `pg` |
| `cmd/spreads` | `spreads`    | `lib/stream` `lib/db` | `kafka` `pg` |
| `cmd/alpha`   | `alpha`      | `lib/stream` `lib/db` | `kafka` `pg` |
| `cmd/api`     |              | `lib/stream` `lib/db` | `kafka` `pg` |

