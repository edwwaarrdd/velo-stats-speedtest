# Velo-stats backend speed test - production stack

Generated 2026-09-10 13:26:04 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

Each backend was started on its own, warmed up, measured, then torn down before
the next one started. No two backends ran at the same time.

This is the **production** profile: each backend's `docker-compose.prod.yml`,
running a real application server with no bind mounts and the already-seeded
SQLite database baked into the image. Every backend is capped at the same four
CPUs, so the runtimes are compared on an equal share of the machine.

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
| 1 | golang | Go 1.25 | net/http, baked-in SQLite, GOMAXPROCS=4 | 0.56 ms | 23656.5 | 0 |
| 2 | nodejs | NestJS 11 / Node 22 | 4 clustered Node processes, baked-in SQLite | 0.84 ms | 20964.2 | 0 |
| 3 | symfony | Symfony 8.1 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 1.12 ms | 16780.6 | 0 |
| 4 | python | Django 6 / gunicorn | gunicorn, 4 preloaded workers, baked-in SQLite | 1.39 ms | 19276.0 | 0 |
| 5 | laravel | Laravel 13 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 4.44 ms | 7018.5 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.28 ms | 0.35 ms | 1.28 ms | 0.62 ms | 0.61 ms |
| `/rides` | 3.27 ms | 4.42 ms | 16.12 ms | 5.78 ms | 5.87 ms |
| `/rides/summary` | 0.56 ms | 0.43 ms | 1.78 ms | 0.97 ms | 1.70 ms |
| `/rides/cost` | 0.46 ms | 0.84 ms | 10.69 ms | 1.12 ms | 1.14 ms |
| `/stations` | 1.11 ms | 1.68 ms | 4.44 ms | 2.52 ms | 1.39 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 15240.6 | 6288.2 | 3473.7 | 6854.9 | 7496.3 |
| `/rides` | 491.4 | 431.9 | 229.6 | 657.4 | 707.4 |
| `/rides/summary` | 2968.6 | 7852.1 | 2252.6 | 4151.4 | 2869.7 |
| `/rides/cost` | 3638.1 | 4216.3 | 347.2 | 3583.6 | 4311.3 |
| `/stations` | 1317.8 | 2175.8 | 715.4 | 1533.3 | 3891.3 |

## Detail per backend

### golang - Go 1.25

net/http, baked-in SQLite, GOMAXPROCS=4 (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.19 ms | 0.28 ms | 0.54 ms | 0.95 ms | 0.95 ms | 0 | 15240.6 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.88 ms | 3.27 ms | 3.72 ms | 4.17 ms | 4.17 ms | 0 | 491.4 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.50 ms | 0.56 ms | 0.63 ms | 0.65 ms | 0.65 ms | 0 | 2968.6 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.41 ms | 0.46 ms | 0.57 ms | 0.81 ms | 0.81 ms | 0 | 3638.1 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.99 ms | 1.11 ms | 1.38 ms | 1.74 ms | 1.74 ms | 0 | 1317.8 |

### nodejs - NestJS 11 / Node 22

4 clustered Node processes, baked-in SQLite (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.28 ms | 0.35 ms | 0.48 ms | 2.88 ms | 2.88 ms | 0 | 6288.2 |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.69 ms | 4.42 ms | 6.05 ms | 8.38 ms | 8.38 ms | 0 | 431.9 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.38 ms | 0.43 ms | 0.52 ms | 0.57 ms | 0.57 ms | 0 | 7852.1 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.73 ms | 0.84 ms | 1.13 ms | 1.52 ms | 1.52 ms | 0 | 4216.3 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.50 ms | 1.68 ms | 1.90 ms | 6.09 ms | 6.09 ms | 0 | 2175.8 |

### laravel - Laravel 13 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 1.05 ms | 1.28 ms | 1.68 ms | 2.29 ms | 2.29 ms | 0 | 3473.7 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 15.48 ms | 16.12 ms | 17.25 ms | 19.11 ms | 19.11 ms | 0 | 229.6 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 1.58 ms | 1.78 ms | 2.38 ms | 2.87 ms | 2.87 ms | 0 | 2252.6 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 9.90 ms | 10.69 ms | 16.71 ms | 22.02 ms | 22.02 ms | 0 | 347.2 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 3.96 ms | 4.44 ms | 5.67 ms | 6.74 ms | 6.74 ms | 0 | 715.4 |

### symfony - Symfony 8.1 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.54 ms | 0.62 ms | 0.77 ms | 1.06 ms | 1.06 ms | 0 | 6854.9 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 5.38 ms | 5.78 ms | 6.32 ms | 6.92 ms | 6.92 ms | 0 | 657.4 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.90 ms | 0.97 ms | 1.06 ms | 1.27 ms | 1.27 ms | 0 | 4151.4 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 1.03 ms | 1.12 ms | 1.21 ms | 1.44 ms | 1.44 ms | 0 | 3583.6 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 2.38 ms | 2.52 ms | 2.61 ms | 3.14 ms | 3.14 ms | 0 | 1533.3 |

### python - Django 6 / gunicorn

gunicorn, 4 preloaded workers, baked-in SQLite (`docker-compose.prod.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 17 B | 0.50 ms | 0.61 ms | 0.75 ms | 0.93 ms | 0.93 ms | 0 | 7496.3 |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.44 ms | 5.87 ms | 6.92 ms | 7.47 ms | 7.47 ms | 0 | 707.4 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.39 ms | 1.70 ms | 1.99 ms | 2.13 ms | 2.13 ms | 0 | 2869.7 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 0.94 ms | 1.14 ms | 1.45 ms | 1.91 ms | 1.91 ms | 0 | 4311.3 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.21 ms | 1.39 ms | 1.66 ms | 1.77 ms | 1.77 ms | 0 | 3891.3 |

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
- **Every backend gets four CPUs**, and the single-threaded runtimes are
  configured with four workers to match: a PHP-FPM pool of four, four gunicorn
  workers, four clustered Node processes, and `GOMAXPROCS=4` for Go. Without
  that cap Go would simply take all twelve cores of the host.
- **Nothing is bind-mounted.** The seeded SQLite database is copied into the
  image at build time, so no request touches the host filesystem.
- **This is still a single container per backend on a laptop.** It is not a
  tuned deployment, and it says nothing about how these stacks behave behind a
  load balancer, with a networked database, or under sustained traffic.

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
