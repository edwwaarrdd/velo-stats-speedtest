# Velo-stats backend speed test - development stack

Generated 2026-09-10 13:38:18 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

Each backend was started on its own, warmed up, measured, then torn down before
the next one started. No two backends ran at the same time.

This is the **development** profile: each backend's own `docker-compose.yml`,
the stack a contributor runs day to day. Source and database are bind-mounted
from the host and no CPU limits are applied.

| Setting | Value |
| --- | --- |
| Requests per endpoint per pass | 500 |
| Concurrency in the concurrent pass | 10 |
| Discarded warmup requests per endpoint | 20 |
| Endpoints measured | 4 |

## Ranking

Median latency is the median of each backend's four per-endpoint medians in the
sequential pass. Throughput is the sum of the four endpoints' concurrent
requests per second. Fastest first.

| # | Backend | Stack | Setup | Median latency | Throughput (req/s) | Errors |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | golang | Go 1.25 | net/http, bind-mounted SQLite, all cores | 1.06 ms | 8559.4 | 0 |
| 2 | nodejs | NestJS 11 / Node 22 | single Node process, bind-mounted source and SQLite | 1.56 ms | 5880.9 | 0 |
| 3 | python | Django 6 / gunicorn | gunicorn, 1 worker, --reload, bind-mounted source | 1.93 ms | 2305.3 | 0 |
| 4 | symfony | Symfony 8.1 / PHP 8.5 | php -S, one request at a time, no OPcache | 5.83 ms | 846.3 | 0 |
| 5 | laravel | Laravel 13 / PHP 8.5 | php artisan serve, one request at a time, no OPcache | 14.89 ms | 511.0 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/rides` | 3.00 ms | 3.98 ms | 20.46 ms | 12.04 ms | 6.46 ms |
| `/rides/summary` | 0.53 ms | 0.40 ms | 4.07 ms | 3.58 ms | 1.93 ms |
| `/rides/cost` | 0.46 ms | 0.79 ms | 14.89 ms | 3.79 ms | 1.42 ms |
| `/stations` | 1.06 ms | 1.56 ms | 6.83 ms | 5.83 ms | 1.63 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/rides` | 508.2 | 264.5 | 49.0 | 89.8 | 163.7 |
| `/rides/summary` | 3051.8 | 3413.1 | 241.6 | 305.0 | 585.8 |
| `/rides/cost` | 3728.4 | 1486.1 | 68.4 | 264.2 | 824.5 |
| `/stations` | 1270.9 | 717.2 | 152.0 | 187.2 | 731.4 |

## Detail per backend

### golang - Go 1.25

net/http, bind-mounted SQLite, all cores (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.75 ms | 3.00 ms | 3.29 ms | 3.45 ms | 3.74 ms | 0 | 508.2 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.46 ms | 0.53 ms | 0.61 ms | 0.65 ms | 0.72 ms | 0 | 3051.8 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.39 ms | 0.46 ms | 0.59 ms | 0.78 ms | 0.90 ms | 0 | 3728.4 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.91 ms | 1.06 ms | 1.35 ms | 1.42 ms | 3.61 ms | 0 | 1270.9 |

### nodejs - NestJS 11 / Node 22

single Node process, bind-mounted source and SQLite (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.70 ms | 3.98 ms | 4.39 ms | 4.80 ms | 5.91 ms | 0 | 264.5 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.36 ms | 0.40 ms | 0.46 ms | 0.59 ms | 2.07 ms | 0 | 3413.1 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.71 ms | 0.79 ms | 0.93 ms | 1.07 ms | 2.49 ms | 0 | 1486.1 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.43 ms | 1.56 ms | 1.70 ms | 1.75 ms | 1.88 ms | 0 | 717.2 |

### laravel - Laravel 13 / PHP 8.5

php artisan serve, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 19.96 ms | 20.46 ms | 22.95 ms | 27.21 ms | 226.59 ms | 0 | 49.0 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.82 ms | 4.07 ms | 4.69 ms | 12.54 ms | 18.65 ms | 0 | 241.6 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 14.43 ms | 14.89 ms | 15.99 ms | 19.75 ms | 34.57 ms | 0 | 68.4 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 6.45 ms | 6.83 ms | 7.55 ms | 11.72 ms | 12.44 ms | 0 | 152.0 |

### symfony - Symfony 8.1 / PHP 8.5

php -S, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 10.98 ms | 12.04 ms | 14.83 ms | 23.30 ms | 25.91 ms | 0 | 89.8 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.24 ms | 3.58 ms | 4.13 ms | 8.25 ms | 13.67 ms | 0 | 305.0 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 3.43 ms | 3.79 ms | 4.22 ms | 4.62 ms | 6.84 ms | 0 | 264.2 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 5.33 ms | 5.83 ms | 7.09 ms | 14.70 ms | 25.47 ms | 0 | 187.2 |

### python - Django 6 / gunicorn

gunicorn, 1 worker, --reload, bind-mounted source (`docker-compose.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.89 ms | 6.46 ms | 8.29 ms | 14.55 ms | 17.64 ms | 0 | 163.7 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.73 ms | 1.93 ms | 2.15 ms | 2.43 ms | 3.25 ms | 0 | 585.8 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 1.26 ms | 1.42 ms | 2.51 ms | 3.79 ms | 15.38 ms | 0 | 824.5 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.48 ms | 1.63 ms | 1.81 ms | 2.02 ms | 3.78 ms | 0 | 731.4 |

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
the part that does. A trivial endpoint that touches neither the database nor an
entity - this harness used to include one, `/_healthcheck`, before dropping
it as measuring nothing the applications do - isolates the cost of accepting a
request and routing it, and separates the backends by well under a millisecond.
Every endpoint below adds work proportional to the rows it serves: 158 rides,
321 stations. What separates the backends by multiples is that per-row cost,
and per-row cost is decided by how much of an object each row is turned into on
the way out. Timing that work inside the Laravel implementation, on 158 rides:

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
