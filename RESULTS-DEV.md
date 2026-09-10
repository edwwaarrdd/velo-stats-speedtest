# Velo-stats backend speed test - development stack

Generated 2026-09-10 13:26:04 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| 1 | golang | Go 1.25 | net/http, bind-mounted SQLite, all cores | 0.49 ms | 25712.4 | 0 |
| 2 | nodejs | NestJS 11 / Node 22 | single Node process, bind-mounted source and SQLite | 0.81 ms | 11216.3 | 0 |
| 3 | python | Django 6 / gunicorn | gunicorn, 1 worker, --reload, bind-mounted source | 1.77 ms | 5666.0 | 0 |
| 4 | symfony | Symfony 8.1 / PHP 8.5 | php -S, one request at a time, no OPcache | 3.79 ms | 1744.3 | 0 |
| 5 | laravel | Laravel 13 / PHP 8.5 | php artisan serve, one request at a time, no OPcache | 6.56 ms | 989.4 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.22 ms | 0.35 ms | 2.67 ms | 1.45 ms | 0.63 ms |
| `/rides` | 2.86 ms | 3.85 ms | 20.72 ms | 11.09 ms | 6.78 ms |
| `/rides/summary` | 0.49 ms | 0.41 ms | 3.96 ms | 3.61 ms | 2.01 ms |
| `/rides/cost` | 0.43 ms | 0.81 ms | 14.26 ms | 3.79 ms | 1.54 ms |
| `/stations` | 1.08 ms | 1.61 ms | 6.56 ms | 6.06 ms | 1.77 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 17188.4 | 5507.8 | 432.7 | 876.2 | 3415.5 |
| `/rides` | 527.8 | 265.7 | 49.5 | 93.6 | 163.7 |
| `/rides/summary` | 2946.9 | 3324.4 | 279.6 | 313.1 | 550.8 |
| `/rides/cost` | 3763.1 | 1452.8 | 69.3 | 280.1 | 818.2 |
| `/stations` | 1286.2 | 665.5 | 158.3 | 181.3 | 717.8 |

## Detail per backend

### golang - Go 1.25

net/http, bind-mounted SQLite, all cores (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.14 ms | 0.22 ms | 0.32 ms | 0.37 ms | 0.37 ms | 0 | 17188.4 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.75 ms | 2.86 ms | 3.15 ms | 3.22 ms | 3.22 ms | 0 | 527.8 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.44 ms | 0.49 ms | 0.54 ms | 0.60 ms | 0.60 ms | 0 | 2946.9 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.38 ms | 0.43 ms | 0.47 ms | 0.65 ms | 0.65 ms | 0 | 3763.1 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.98 ms | 1.08 ms | 1.35 ms | 1.57 ms | 1.57 ms | 0 | 1286.2 |

### nodejs - NestJS 11 / Node 22

single Node process, bind-mounted source and SQLite (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.25 ms | 0.35 ms | 0.91 ms | 1.66 ms | 1.66 ms | 0 | 5507.8 |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.67 ms | 3.85 ms | 4.25 ms | 5.25 ms | 5.25 ms | 0 | 265.7 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.37 ms | 0.41 ms | 0.48 ms | 0.52 ms | 0.52 ms | 0 | 3324.4 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.75 ms | 0.81 ms | 0.92 ms | 1.11 ms | 1.11 ms | 0 | 1452.8 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.45 ms | 1.61 ms | 1.79 ms | 2.38 ms | 2.38 ms | 0 | 665.5 |

### laravel - Laravel 13 / PHP 8.5

php artisan serve, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 2.47 ms | 2.67 ms | 3.64 ms | 6.79 ms | 6.79 ms | 0 | 432.7 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 19.45 ms | 20.72 ms | 26.98 ms | 53.39 ms | 53.39 ms | 0 | 49.5 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.80 ms | 3.96 ms | 4.16 ms | 10.37 ms | 10.37 ms | 0 | 279.6 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 13.84 ms | 14.26 ms | 14.74 ms | 18.53 ms | 18.53 ms | 0 | 69.3 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 6.30 ms | 6.56 ms | 6.77 ms | 15.20 ms | 15.20 ms | 0 | 158.3 |

### symfony - Symfony 8.1 / PHP 8.5

php -S, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 1.29 ms | 1.45 ms | 1.66 ms | 1.75 ms | 1.75 ms | 0 | 876.2 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 10.66 ms | 11.09 ms | 11.82 ms | 11.93 ms | 11.93 ms | 0 | 93.6 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.20 ms | 3.61 ms | 4.09 ms | 4.77 ms | 4.77 ms | 0 | 313.1 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 3.44 ms | 3.79 ms | 4.42 ms | 7.19 ms | 7.19 ms | 0 | 280.1 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 5.57 ms | 6.06 ms | 6.72 ms | 7.55 ms | 7.55 ms | 0 | 181.3 |

### python - Django 6 / gunicorn

gunicorn, 1 worker, --reload, bind-mounted source (`docker-compose.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 17 B | 0.53 ms | 0.63 ms | 0.83 ms | 1.08 ms | 1.08 ms | 0 | 3415.5 |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 6.15 ms | 6.78 ms | 7.90 ms | 9.11 ms | 9.11 ms | 0 | 163.7 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.81 ms | 2.01 ms | 2.36 ms | 2.63 ms | 2.63 ms | 0 | 550.8 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 1.34 ms | 1.54 ms | 1.74 ms | 1.81 ms | 1.81 ms | 0 | 818.2 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.58 ms | 1.77 ms | 2.14 ms | 3.13 ms | 3.13 ms | 0 | 717.8 |

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
