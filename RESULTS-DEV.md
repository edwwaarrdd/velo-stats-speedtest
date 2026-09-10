# Velo-stats backend speed test - development stack

Generated 2026-09-10 11:15:41 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| 1 | golang | Go 1.25 | net/http, bind-mounted SQLite, all cores | 0.56 ms | 23680.6 | 0 |
| 2 | nodejs | NestJS 11 / Node 22 | single Node process, bind-mounted source and SQLite | 0.86 ms | 10046.1 | 0 |
| 3 | python | Django 6 / gunicorn | gunicorn, 1 worker, --reload, bind-mounted source | 1.68 ms | 4783.3 | 0 |
| 4 | symfony | Symfony 8.1 / PHP 8.5 | php -S, one request at a time, no OPcache | 4.24 ms | 1657.5 | 0 |
| 5 | php | Laravel 13 / PHP 8.4 | php artisan serve, one request at a time, no OPcache | 7.47 ms | 860.7 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | php | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.23 ms | 0.38 ms | 2.84 ms | 1.53 ms | 0.71 ms |
| `/rides` | 2.94 ms | 4.07 ms | 21.13 ms | 12.00 ms | 8.29 ms |
| `/rides/summary` | 0.56 ms | 0.43 ms | 4.79 ms | 4.01 ms | 2.07 ms |
| `/rides/cost` | 0.43 ms | 0.86 ms | 15.97 ms | 4.24 ms | 1.52 ms |
| `/stations` | 1.05 ms | 1.73 ms | 7.47 ms | 6.07 ms | 1.68 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | php | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 15254.4 | 4513.8 | 370.0 | 838.7 | 2508.2 |
| `/rides` | 506.2 | 257.6 | 48.1 | 87.6 | 130.5 |
| `/rides/summary` | 2938.4 | 3182.7 | 236.2 | 289.7 | 548.3 |
| `/rides/cost` | 3729.4 | 1450.3 | 61.9 | 265.2 | 848.3 |
| `/stations` | 1252.2 | 641.6 | 144.5 | 176.4 | 747.9 |

## Detail per backend

### golang - Go 1.25

net/http, bind-mounted SQLite, all cores (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.16 ms | 0.23 ms | 0.40 ms | 0.43 ms | 0.43 ms | 0 | 15254.4 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.72 ms | 2.94 ms | 3.31 ms | 4.85 ms | 4.85 ms | 0 | 506.2 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.48 ms | 0.56 ms | 0.63 ms | 0.69 ms | 0.69 ms | 0 | 2938.4 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.39 ms | 0.43 ms | 0.69 ms | 0.76 ms | 0.76 ms | 0 | 3729.4 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.97 ms | 1.05 ms | 1.34 ms | 1.39 ms | 1.39 ms | 0 | 1252.2 |

### nodejs - NestJS 11 / Node 22

single Node process, bind-mounted source and SQLite (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.29 ms | 0.38 ms | 0.52 ms | 2.94 ms | 2.94 ms | 0 | 4513.8 |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.80 ms | 4.07 ms | 4.52 ms | 5.02 ms | 5.02 ms | 0 | 257.6 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.36 ms | 0.43 ms | 0.50 ms | 0.55 ms | 0.55 ms | 0 | 3182.7 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.73 ms | 0.86 ms | 1.02 ms | 1.19 ms | 1.19 ms | 0 | 1450.3 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.55 ms | 1.73 ms | 1.98 ms | 3.47 ms | 3.47 ms | 0 | 641.6 |

### php - Laravel 13 / PHP 8.4

php artisan serve, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 2.63 ms | 2.84 ms | 3.08 ms | 3.34 ms | 3.34 ms | 0 | 370.0 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 20.12 ms | 21.13 ms | 23.54 ms | 24.86 ms | 24.86 ms | 0 | 48.1 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 4.24 ms | 4.79 ms | 5.69 ms | 7.69 ms | 7.69 ms | 0 | 236.2 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 15.17 ms | 15.97 ms | 21.40 ms | 33.85 ms | 33.85 ms | 0 | 61.9 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 6.85 ms | 7.47 ms | 9.15 ms | 11.49 ms | 11.49 ms | 0 | 144.5 |

### symfony - Symfony 8.1 / PHP 8.5

php -S, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 1.39 ms | 1.53 ms | 1.88 ms | 2.05 ms | 2.05 ms | 0 | 838.7 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 11.34 ms | 12.00 ms | 13.08 ms | 17.01 ms | 17.01 ms | 0 | 87.6 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.49 ms | 4.01 ms | 4.58 ms | 8.03 ms | 8.03 ms | 0 | 289.7 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 3.71 ms | 4.24 ms | 4.81 ms | 4.94 ms | 4.94 ms | 0 | 265.2 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 5.68 ms | 6.07 ms | 6.46 ms | 6.72 ms | 6.72 ms | 0 | 176.4 |

### python - Django 6 / gunicorn

gunicorn, 1 worker, --reload, bind-mounted source (`docker-compose.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 17 B | 0.52 ms | 0.71 ms | 1.03 ms | 2.33 ms | 2.33 ms | 0 | 2508.2 |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 7.10 ms | 8.29 ms | 12.11 ms | 23.78 ms | 23.78 ms | 0 | 130.5 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.75 ms | 2.07 ms | 2.55 ms | 2.76 ms | 2.76 ms | 0 | 548.3 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 1.25 ms | 1.52 ms | 2.09 ms | 2.33 ms | 2.33 ms | 0 | 848.3 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.49 ms | 1.68 ms | 1.99 ms | 2.43 ms | 2.43 ms | 0 | 747.9 |

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

## What this measures, and what it does not

A ranking of five backends invites being read as a ranking of five frameworks.
It is not one, and the numbers themselves say so.

Split each backend's latency into the part that does not depend on the data and
the part that does. `/_healthcheck` returns a fixed string and touches
neither the database nor an entity, so it isolates the cost of accepting a
request and routing it. Every other endpoint adds work proportional to the rows
it serves: 158 rides, 321 stations.

Measured that way, the fixed cost separates the backends by well under a
millisecond. What separates them by multiples is the per-row cost, and per-row
cost is decided by how much of an object each row is turned into on the way out.
Timing that work inside the Laravel implementation, on 158 rides:

| Stage | Time |
| --- | --- |
| The database answering the query, rows as arrays | 0.27 ms |
| Hydrating those rows into 158 models, weather included | 2.80 ms |
| Mapping each model through an API resource | 9.64 ms |

Roughly 97% of the data work on that endpoint is constructing PHP objects, not
querying. The cost endpoint is the same story in miniature: reading one
timestamp column as date objects takes 1.87 ms, and reading the same column as
plain strings takes 0.04 ms.

So a backend that hydrates entities and serialises them will lose to one that
maps arrays, in any language, by a margin that grows with the row count and has
almost nothing to do with the framework wrapped around it. All five backends
here hydrate their rows, which is what makes them comparable. Change any one of
them to map arrays instead and it will jump the ranking, without its framework
having changed at all.

Read these tables as a comparison of five implementations, then. Not of Go
against PHP, and not of one framework against another.
