# Velo-stats backend speed test - production stack

Generated 2026-09-10 13:28:59 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

Each backend was started on its own, warmed up, measured, then torn down before
the next one started. No two backends ran at the same time.

This is the **production** profile: each backend's `docker-compose.prod.yml`,
running a real application server with no bind mounts and the already-seeded
SQLite database baked into the image. Every backend is capped at the same four
CPUs, so the runtimes are compared on an equal share of the machine.

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
| 1 | golang | Go 1.25 | net/http, baked-in SQLite, GOMAXPROCS=4 | 1.11 ms | 8444.0 | 0 |
| 2 | nodejs | NestJS 11 / Node 22 | 4 clustered Node processes, baked-in SQLite | 1.56 ms | 20481.4 | 0 |
| 3 | python | Django 6 / gunicorn | gunicorn, 4 preloaded workers, baked-in SQLite | 1.86 ms | 11092.2 | 0 |
| 4 | symfony | Symfony 8.1 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 2.69 ms | 9263.0 | 0 |
| 5 | laravel | Laravel 13 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 9.85 ms | 4045.1 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/rides` | 3.33 ms | 3.80 ms | 15.14 ms | 6.27 ms | 6.20 ms |
| `/rides/summary` | 0.51 ms | 0.38 ms | 1.59 ms | 1.20 ms | 1.86 ms |
| `/rides/cost` | 0.45 ms | 0.78 ms | 9.85 ms | 1.33 ms | 1.18 ms |
| `/stations` | 1.11 ms | 1.56 ms | 4.39 ms | 2.69 ms | 1.38 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/rides` | 470.1 | 790.8 | 244.7 | 631.0 | 553.4 |
| `/rides/summary` | 3135.0 | 12354.3 | 2555.7 | 4016.6 | 2831.3 |
| `/rides/cost` | 3680.3 | 5205.8 | 368.3 | 3216.5 | 4295.8 |
| `/stations` | 1158.6 | 2130.5 | 876.4 | 1399.0 | 3411.7 |

## Detail per backend

### golang - Go 1.25

net/http, baked-in SQLite, GOMAXPROCS=4 (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.94 ms | 3.33 ms | 3.98 ms | 4.76 ms | 6.37 ms | 0 | 470.1 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.43 ms | 0.51 ms | 0.61 ms | 0.85 ms | 1.16 ms | 0 | 3135.0 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.38 ms | 0.45 ms | 0.65 ms | 0.89 ms | 1.51 ms | 0 | 3680.3 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.97 ms | 1.11 ms | 1.46 ms | 1.73 ms | 2.32 ms | 0 | 1158.6 |

### nodejs - NestJS 11 / Node 22

4 clustered Node processes, baked-in SQLite (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.61 ms | 3.80 ms | 4.12 ms | 4.34 ms | 4.75 ms | 0 | 790.8 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.30 ms | 0.38 ms | 0.44 ms | 0.51 ms | 0.99 ms | 0 | 12354.3 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.68 ms | 0.78 ms | 0.98 ms | 1.20 ms | 2.32 ms | 0 | 5205.8 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.45 ms | 1.56 ms | 2.05 ms | 2.81 ms | 6.58 ms | 0 | 2130.5 |

### laravel - Laravel 13 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 14.61 ms | 15.14 ms | 16.52 ms | 18.84 ms | 29.89 ms | 0 | 244.7 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 1.45 ms | 1.59 ms | 1.81 ms | 4.09 ms | 12.14 ms | 0 | 2555.7 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 9.39 ms | 9.85 ms | 11.73 ms | 18.56 ms | 29.66 ms | 0 | 368.3 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 4.00 ms | 4.39 ms | 5.43 ms | 5.81 ms | 6.05 ms | 0 | 876.4 |

### symfony - Symfony 8.1 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 5.55 ms | 6.27 ms | 7.58 ms | 8.93 ms | 9.88 ms | 0 | 631.0 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.90 ms | 1.20 ms | 1.66 ms | 2.09 ms | 2.24 ms | 0 | 4016.6 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 1.04 ms | 1.33 ms | 1.82 ms | 2.22 ms | 3.45 ms | 0 | 3216.5 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 2.43 ms | 2.69 ms | 3.32 ms | 4.17 ms | 5.03 ms | 0 | 1399.0 |

### python - Django 6 / gunicorn

gunicorn, 4 preloaded workers, baked-in SQLite (`docker-compose.prod.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.47 ms | 6.20 ms | 7.00 ms | 8.27 ms | 21.12 ms | 0 | 553.4 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.50 ms | 1.86 ms | 2.31 ms | 2.88 ms | 3.42 ms | 0 | 2831.3 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 0.95 ms | 1.18 ms | 1.51 ms | 1.80 ms | 4.04 ms | 0 | 4295.8 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.22 ms | 1.38 ms | 1.59 ms | 1.71 ms | 2.75 ms | 0 | 3411.7 |

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
