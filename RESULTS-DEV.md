# Velo-stats backend speed test - development stack

Generated 2026-09-10 15:47:29 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

Each backend was started on its own, warmed up, measured, then torn down before
the next one started. No two backends ran at the same time.

**Every figure below is based on 1000 requests per endpoint,
in each of two passes.** A further 2000 requests per endpoint
are issued first and discarded, which is set by the slowest runtime to warm
rather than the fastest: the JVM needs thousands of calls to finish compiling,
where the others are at full speed within tens.

This is the **development** profile: each backend's own `docker-compose.yml`,
the stack a contributor runs day to day. Source and database are bind-mounted
from the host and no CPU limits are applied.

| Setting | Value |
| --- | --- |
| Requests per endpoint per pass | 1000 |
| Concurrency in the concurrent pass | 10 |
| Discarded warmup requests per endpoint | 2000 |
| Endpoints measured | 4 |

## Ranking

Median latency is the median of each backend's four per-endpoint medians in the
sequential pass. Throughput is the sum of the four endpoints' concurrent
requests per second. Fastest first.

| # | Backend | Stack | Setup | Median latency | Throughput (req/s) | Errors |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | golang | Go 1.25 | net/http, bind-mounted SQLite, all cores | 1.12 ms | 8287.4 | 0 |
| 2 | java | Spring Boot 4.1 / Java 25 | embedded server, bind-mounted SQLite, all cores | 1.15 ms | 32379.6 | 0 |
| 3 | nodejs | NestJS 11 / Node 22 | single Node process, bind-mounted source and SQLite | 1.58 ms | 6181.6 | 0 |
| 4 | python | Django 6 / gunicorn | gunicorn, 1 worker, --reload, bind-mounted source | 2.02 ms | 2212.1 | 0 |
| 5 | symfony | Symfony 8.1 / PHP 8.5 | php -S, one request at a time, no OPcache | 5.64 ms | 845.3 | 0 |
| 6 | laravel | Laravel 13 / PHP 8.5 | php artisan serve, one request at a time, no OPcache | 15.06 ms | 510.9 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | java | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- | --- |
| `/rides` | 2.93 ms | 2.48 ms | 3.91 ms | 20.84 ms | 11.53 ms | 6.35 ms |
| `/rides/summary` | 0.52 ms | 0.33 ms | 0.35 ms | 4.18 ms | 3.95 ms | 2.02 ms |
| `/rides/cost` | 0.42 ms | 0.31 ms | 0.74 ms | 15.06 ms | 3.88 ms | 1.48 ms |
| `/stations` | 1.12 ms | 1.15 ms | 1.58 ms | 8.28 ms | 5.64 ms | 1.72 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | java | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- | --- |
| `/rides` | 504.5 | 1285.1 | 260.4 | 47.5 | 81.2 | 165.2 |
| `/rides/summary` | 2706.0 | 12977.2 | 3670.2 | 258.9 | 296.1 | 578.9 |
| `/rides/cost` | 3765.4 | 14152.9 | 1566.2 | 66.9 | 283.2 | 734.8 |
| `/stations` | 1311.6 | 3964.3 | 684.8 | 137.7 | 184.7 | 733.2 |

## Detail per backend

### golang - Go 1.25

net/http, bind-mounted SQLite, all cores (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.67 ms | 2.93 ms | 3.23 ms | 3.42 ms | 3.79 ms | 0 | 504.5 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.45 ms | 0.52 ms | 0.57 ms | 0.62 ms | 0.77 ms | 0 | 2706.0 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.37 ms | 0.42 ms | 0.52 ms | 0.72 ms | 0.82 ms | 0 | 3765.4 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.96 ms | 1.12 ms | 1.37 ms | 1.49 ms | 2.35 ms | 0 | 1311.6 |

### java - Spring Boot 4.1 / Java 25

embedded server, bind-mounted SQLite, all cores (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.18 ms | 2.48 ms | 3.48 ms | 4.52 ms | 8.40 ms | 0 | 1285.1 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.26 ms | 0.33 ms | 0.43 ms | 0.86 ms | 1.42 ms | 0 | 12977.2 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.25 ms | 0.31 ms | 0.37 ms | 0.84 ms | 1.57 ms | 0 | 14152.9 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.04 ms | 1.15 ms | 1.27 ms | 1.79 ms | 2.43 ms | 0 | 3964.3 |

### nodejs - NestJS 11 / Node 22

single Node process, bind-mounted source and SQLite (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.72 ms | 3.91 ms | 4.24 ms | 5.43 ms | 10.79 ms | 0 | 260.4 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.30 ms | 0.35 ms | 0.40 ms | 0.46 ms | 1.40 ms | 0 | 3670.2 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.65 ms | 0.74 ms | 0.84 ms | 0.99 ms | 2.22 ms | 0 | 1566.2 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.41 ms | 1.58 ms | 1.85 ms | 2.21 ms | 4.25 ms | 0 | 684.8 |

### laravel - Laravel 13 / PHP 8.5

php artisan serve, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 20.33 ms | 20.84 ms | 22.33 ms | 32.51 ms | 234.86 ms | 0 | 47.5 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.82 ms | 4.18 ms | 5.02 ms | 8.31 ms | 14.54 ms | 0 | 258.9 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 14.63 ms | 15.06 ms | 16.02 ms | 19.01 ms | 35.66 ms | 0 | 66.9 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 7.00 ms | 8.28 ms | 9.37 ms | 12.44 ms | 27.14 ms | 0 | 137.7 |

### symfony - Symfony 8.1 / PHP 8.5

php -S, one request at a time, no OPcache (`docker-compose.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 10.70 ms | 11.53 ms | 13.82 ms | 21.07 ms | 140.97 ms | 0 | 81.2 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 3.37 ms | 3.95 ms | 4.83 ms | 8.63 ms | 17.56 ms | 0 | 296.1 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 3.48 ms | 3.88 ms | 4.37 ms | 4.74 ms | 7.55 ms | 0 | 283.2 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 5.24 ms | 5.64 ms | 6.12 ms | 6.93 ms | 9.25 ms | 0 | 184.7 |

### python - Django 6 / gunicorn

gunicorn, 1 worker, --reload, bind-mounted source (`docker-compose.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.85 ms | 6.35 ms | 7.23 ms | 9.41 ms | 16.38 ms | 0 | 165.2 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.78 ms | 2.02 ms | 2.33 ms | 2.65 ms | 7.22 ms | 0 | 578.9 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 1.25 ms | 1.48 ms | 1.62 ms | 1.76 ms | 5.01 ms | 0 | 734.8 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.51 ms | 1.72 ms | 1.96 ms | 4.52 ms | 11.87 ms | 0 | 733.2 |

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

A ranking of six backends invites being read as a ranking of six frameworks.
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
almost nothing to do with the framework wrapped around it. All six backends
here hydrate their rows, which is what makes them comparable. Change any one of
them to map arrays instead and it will jump the ranking, without its framework
having changed at all.

Read these tables as a comparison of six implementations, then. Not of Go
against PHP, and not of one framework against another.
