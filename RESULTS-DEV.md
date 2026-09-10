# Velo-stats backend speed test - development stack

Generated 2026-09-10 10:46:16 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

Each backend was started on its own, warmed up, measured, then torn down before
the next one started. No two backends ran at the same time.

This is the **development** profile: each backend's own `docker-compose.yml`,
the stack a contributor runs day to day. Source and database are bind-mounted
from the host and no CPU limits are applied.

| Setting | Value |
| --- | --- |
| Requests per endpoint per pass | 100 |
| Concurrency in the concurrent pass | 10 |
| Discarded warmup requests per endpoint | 20 |
| Endpoints measured | 5 |

## Ranking

Median latency is the median of each backend's five per-endpoint medians in the
sequential pass. Throughput is the sum of the five endpoints' concurrent
requests per second. Fastest first.

| # | Backend | Stack | Setup | Median latency | Throughput (req/s) | Errors |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | golang | Go 1.25 | net/http, bind-mounted SQLite, all cores | 0.54 ms | 22859.8 | 0 |
| 2 | nodejs | NestJS 11 / Node 22 | single Node process, bind-mounted source and SQLite | 0.83 ms | 9767.5 | 0 |
| 3 | python | Django 6 / gunicorn | gunicorn, 1 worker, --reload, bind-mounted source | 1.61 ms | 5469.4 | 0 |
| 4 | symfony | Symfony 8.1 / PHP 8.5 | php -S, one request at a time, no OPcache | 3.63 ms | 1891.1 | 0 |
| 5 | php | Laravel 13 / PHP 8.4 | php artisan serve, one request at a time, no OPcache | 7.18 ms | 855.2 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | php | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.29 ms | 0.42 ms | 2.74 ms | 1.68 ms | 0.58 ms |
| `/rides` | 3.17 ms | 4.14 ms | 22.94 ms | 7.53 ms | 6.20 ms |
| `/rides/summary` | 0.54 ms | 0.42 ms | 4.30 ms | 4.13 ms | 1.95 ms |
| `/rides/cost` | 0.45 ms | 0.83 ms | 16.39 ms | 3.49 ms | 1.44 ms |
| `/stations` | 1.09 ms | 1.65 ms | 7.18 ms | 3.63 ms | 1.61 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | php | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 14325.3 | 4419.3 | 324.4 | 829.2 | 3069.8 |
| `/rides` | 487.7 | 252.3 | 46.3 | 165.7 | 167.4 |
| `/rides/summary` | 3089.9 | 3031.9 | 267.8 | 290.8 | 611.5 |
| `/rides/cost` | 3701.0 | 1412.9 | 66.0 | 312.7 | 880.5 |
| `/stations` | 1255.8 | 651.0 | 150.7 | 292.8 | 740.3 |

## Detail per backend

### golang - Go 1.25

net/http, bind-mounted SQLite, all cores (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.18 ms | 0.29 ms | 1.21 ms | 3.88 ms | 3.88 ms | 0 | 14325.3 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.85 ms | 3.17 ms | 3.47 ms | 3.61 ms | 3.61 ms | 0 | 487.7 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.49 ms | 0.54 ms | 0.62 ms | 1.21 ms | 1.21 ms | 0 | 3089.9 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.42 ms | 0.45 ms | 0.66 ms | 0.73 ms | 0.73 ms | 0 | 3701.0 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.00 ms | 1.09 ms | 1.37 ms | 1.54 ms | 1.54 ms | 0 | 1255.8 |

### nodejs - NestJS 11 / Node 22

single Node process, bind-mounted source and SQLite (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.29 ms | 0.42 ms | 0.60 ms | 3.93 ms | 3.93 ms | 0 | 4419.3 |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.81 ms | 4.14 ms | 4.59 ms | 5.20 ms | 5.20 ms | 0 | 252.3 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.37 ms | 0.42 ms | 0.50 ms | 0.71 ms | 0.71 ms | 0 | 3031.9 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.75 ms | 0.83 ms | 0.97 ms | 1.39 ms | 1.39 ms | 0 | 1412.9 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.51 ms | 1.65 ms | 1.79 ms | 1.89 ms | 1.89 ms | 0 | 651.0 |

### php - Laravel 13 / PHP 8.4

php artisan serve, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 2.55 ms | 2.74 ms | 3.53 ms | 3.75 ms | 3.75 ms | 0 | 324.4 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 21.27 ms | 22.94 ms | 26.80 ms | 31.02 ms | 31.02 ms | 0 | 46.3 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.96 ms | 4.30 ms | 4.81 ms | 10.18 ms | 10.18 ms | 0 | 267.8 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 14.47 ms | 16.39 ms | 21.01 ms | 34.41 ms | 34.41 ms | 0 | 66.0 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 6.54 ms | 7.18 ms | 13.35 ms | 20.77 ms | 20.77 ms | 0 | 150.7 |

### symfony - Symfony 8.1 / PHP 8.5

php -S, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 1.48 ms | 1.68 ms | 2.22 ms | 7.43 ms | 7.43 ms | 0 | 829.2 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 5.87 ms | 7.53 ms | 10.47 ms | 27.52 ms | 27.52 ms | 0 | 165.7 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.57 ms | 4.13 ms | 5.31 ms | 6.60 ms | 6.60 ms | 0 | 290.8 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 3.10 ms | 3.49 ms | 4.07 ms | 4.28 ms | 4.28 ms | 0 | 312.7 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 3.30 ms | 3.63 ms | 4.09 ms | 10.60 ms | 10.60 ms | 0 | 292.8 |

### python - Django 6 / gunicorn

gunicorn, 1 worker, --reload, bind-mounted source (`docker-compose.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 17 B | 0.47 ms | 0.58 ms | 0.81 ms | 0.92 ms | 0.92 ms | 0 | 3069.8 |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.79 ms | 6.20 ms | 6.77 ms | 7.07 ms | 7.07 ms | 0 | 167.4 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.78 ms | 1.95 ms | 2.24 ms | 2.44 ms | 2.44 ms | 0 | 611.5 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 1.31 ms | 1.44 ms | 1.58 ms | 1.98 ms | 1.98 ms | 0 | 880.5 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.50 ms | 1.61 ms | 1.86 ms | 2.21 ms | 2.21 ms | 0 | 740.3 |

## How to read this, and what it does not say

- **Compare the body sizes in the detail tables.** If one backend's `/rides`
  response is far smaller than the others, its database was never seeded and its
  numbers are not comparable.
- **Sequential latency and concurrent throughput answer different questions.**
  The first is how fast one request is served on an idle server; the second is
  how much work the server gets through under load. A runtime can win one and
  lose the other.
- **This is one run on one machine.** Latency at this scale is sensitive to
  background load. Re-run before drawing a conclusion from a small gap.
- **These are development stacks, and two of them are deliberately slow.**
  Laravel is served by `php artisan serve`, a single-process development server
  that handles one request at a time with no OPcache. Django runs one gunicorn
  worker with `--reload`. Neither is what production looks like. See
  `RESULTS-PROD.md` for the same measurement against real application servers.
- **Source and database are bind-mounted**, so on macOS every file read goes
  through Docker Desktop's filesystem layer.
