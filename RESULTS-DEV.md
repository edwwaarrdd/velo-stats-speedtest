# Velo-stats backend speed test - development stack

Generated 2026-09-10 09:13:06 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| 1 | golang | Go 1.25 | net/http, bind-mounted SQLite, all cores | 0.54 ms | 24901.9 | 0 |
| 2 | nodejs | NestJS 11 / Node 22 | single Node process, bind-mounted source and SQLite | 0.80 ms | 10222.3 | 0 |
| 3 | python | Django 6 / gunicorn | gunicorn, 1 worker, --reload, bind-mounted source | 1.58 ms | 5443.9 | 0 |
| 4 | php | Laravel 13 / PHP 8.4 | php artisan serve, one request at a time, no OPcache | 7.01 ms | 882.0 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | php | python |
| --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.24 ms | 0.35 ms | 3.02 ms | 0.68 ms |
| `/rides` | 3.02 ms | 3.85 ms | 19.74 ms | 6.11 ms |
| `/rides/summary` | 0.54 ms | 0.42 ms | 3.94 ms | 1.99 ms |
| `/rides/cost` | 0.47 ms | 0.80 ms | 14.85 ms | 1.40 ms |
| `/stations` | 1.07 ms | 1.65 ms | 7.01 ms | 1.58 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | php | python |
| --- | --- | --- | --- | --- |
| `/_healthcheck` | 16739.6 | 4685.5 | 354.3 | 3163.0 |
| `/rides` | 489.6 | 265.7 | 52.1 | 168.5 |
| `/rides/summary` | 2866.1 | 3230.5 | 268.3 | 599.3 |
| `/rides/cost` | 3541.3 | 1426.6 | 67.8 | 798.6 |
| `/stations` | 1265.3 | 614.0 | 139.4 | 714.6 |

## Detail per backend

### golang - Go 1.25

net/http, bind-mounted SQLite, all cores (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.17 ms | 0.24 ms | 0.34 ms | 0.48 ms | 0.48 ms | 0 | 16739.6 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.77 ms | 3.02 ms | 3.35 ms | 3.61 ms | 3.61 ms | 0 | 489.6 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.48 ms | 0.54 ms | 0.60 ms | 0.61 ms | 0.61 ms | 0 | 2866.1 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.40 ms | 0.47 ms | 0.68 ms | 1.25 ms | 1.25 ms | 0 | 3541.3 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.94 ms | 1.07 ms | 1.36 ms | 1.44 ms | 1.44 ms | 0 | 1265.3 |

### nodejs - NestJS 11 / Node 22

single Node process, bind-mounted source and SQLite (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.29 ms | 0.35 ms | 0.55 ms | 0.65 ms | 0.65 ms | 0 | 4685.5 |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.67 ms | 3.85 ms | 4.21 ms | 4.32 ms | 4.32 ms | 0 | 265.7 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.39 ms | 0.42 ms | 0.48 ms | 0.50 ms | 0.50 ms | 0 | 3230.5 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.72 ms | 0.80 ms | 0.90 ms | 1.26 ms | 1.26 ms | 0 | 1426.6 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.54 ms | 1.65 ms | 1.87 ms | 2.07 ms | 2.07 ms | 0 | 614.0 |

### php - Laravel 13 / PHP 8.4

php artisan serve, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 2.65 ms | 3.02 ms | 3.58 ms | 3.81 ms | 3.81 ms | 0 | 354.3 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 19.00 ms | 19.74 ms | 20.99 ms | 21.78 ms | 21.78 ms | 0 | 52.1 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.77 ms | 3.94 ms | 4.10 ms | 8.27 ms | 8.27 ms | 0 | 268.3 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 13.81 ms | 14.85 ms | 18.58 ms | 19.52 ms | 19.52 ms | 0 | 67.8 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 6.46 ms | 7.01 ms | 7.93 ms | 12.70 ms | 12.70 ms | 0 | 139.4 |

### python - Django 6 / gunicorn

gunicorn, 1 worker, --reload, bind-mounted source (`docker-compose.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 17 B | 0.53 ms | 0.68 ms | 0.91 ms | 1.21 ms | 1.21 ms | 0 | 3163.0 |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.74 ms | 6.11 ms | 7.00 ms | 8.96 ms | 8.96 ms | 0 | 168.5 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.79 ms | 1.99 ms | 2.76 ms | 4.54 ms | 4.54 ms | 0 | 599.3 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 1.26 ms | 1.40 ms | 1.61 ms | 2.69 ms | 2.69 ms | 0 | 798.6 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.48 ms | 1.58 ms | 1.70 ms | 1.89 ms | 1.89 ms | 0 | 714.6 |

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
