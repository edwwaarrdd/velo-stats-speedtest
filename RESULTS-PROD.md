# Velo-stats backend speed test - production stack

Generated 2026-09-10 13:38:18 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| 1 | golang | Go 1.25 | net/http, baked-in SQLite, GOMAXPROCS=4 | 1.09 ms | 8690.9 | 0 |
| 2 | python | Django 6 / gunicorn | gunicorn, 4 preloaded workers, baked-in SQLite | 1.64 ms | 11237.0 | 0 |
| 3 | nodejs | NestJS 11 / Node 22 | 4 clustered Node processes, baked-in SQLite | 1.67 ms | 16426.9 | 0 |
| 4 | symfony | Symfony 8.1 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 2.56 ms | 10523.6 | 0 |
| 5 | laravel | Laravel 13 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 9.84 ms | 4199.0 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/rides` | 2.95 ms | 3.88 ms | 15.20 ms | 5.52 ms | 5.74 ms |
| `/rides/summary` | 0.53 ms | 0.41 ms | 1.60 ms | 0.98 ms | 1.64 ms |
| `/rides/cost` | 0.43 ms | 0.77 ms | 9.84 ms | 1.14 ms | 1.07 ms |
| `/stations` | 1.09 ms | 1.67 ms | 4.01 ms | 2.56 ms | 1.41 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/rides` | 521.3 | 751.0 | 250.8 | 643.1 | 646.2 |
| `/rides/summary` | 3115.5 | 10261.7 | 2609.4 | 4410.7 | 2906.4 |
| `/rides/cost` | 3770.6 | 3053.0 | 394.5 | 3875.5 | 4312.4 |
| `/stations` | 1283.4 | 2361.2 | 944.3 | 1594.3 | 3372.0 |

## Detail per backend

### golang - Go 1.25

net/http, baked-in SQLite, GOMAXPROCS=4 (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.75 ms | 2.95 ms | 3.35 ms | 4.17 ms | 4.42 ms | 0 | 521.3 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.47 ms | 0.53 ms | 0.60 ms | 0.67 ms | 0.83 ms | 0 | 3115.5 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.38 ms | 0.43 ms | 0.55 ms | 0.77 ms | 1.06 ms | 0 | 3770.6 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.92 ms | 1.09 ms | 1.33 ms | 1.43 ms | 1.50 ms | 0 | 1283.4 |

### nodejs - NestJS 11 / Node 22

4 clustered Node processes, baked-in SQLite (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.62 ms | 3.88 ms | 4.50 ms | 4.98 ms | 5.57 ms | 0 | 751.0 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.32 ms | 0.41 ms | 0.51 ms | 0.86 ms | 1.33 ms | 0 | 10261.7 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.69 ms | 0.77 ms | 0.95 ms | 1.15 ms | 5.57 ms | 0 | 3053.0 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.48 ms | 1.67 ms | 1.94 ms | 2.23 ms | 9.13 ms | 0 | 2361.2 |

### laravel - Laravel 13 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 14.79 ms | 15.20 ms | 16.25 ms | 17.92 ms | 23.51 ms | 0 | 250.8 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 1.45 ms | 1.60 ms | 1.71 ms | 1.76 ms | 2.00 ms | 0 | 2609.4 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 9.39 ms | 9.84 ms | 10.92 ms | 13.17 ms | 27.68 ms | 0 | 394.5 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 3.84 ms | 4.01 ms | 4.27 ms | 4.90 ms | 5.14 ms | 0 | 944.3 |

### symfony - Symfony 8.1 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 5.19 ms | 5.52 ms | 6.10 ms | 6.68 ms | 7.28 ms | 0 | 643.1 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.86 ms | 0.98 ms | 1.08 ms | 1.18 ms | 1.48 ms | 0 | 4410.7 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 1.04 ms | 1.14 ms | 1.22 ms | 1.27 ms | 1.65 ms | 0 | 3875.5 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 2.40 ms | 2.56 ms | 2.74 ms | 2.97 ms | 3.75 ms | 0 | 1594.3 |

### python - Django 6 / gunicorn

gunicorn, 4 preloaded workers, baked-in SQLite (`docker-compose.prod.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.31 ms | 5.74 ms | 6.35 ms | 7.37 ms | 25.79 ms | 0 | 646.2 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.45 ms | 1.64 ms | 1.87 ms | 2.01 ms | 2.96 ms | 0 | 2906.4 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 0.92 ms | 1.07 ms | 1.26 ms | 1.38 ms | 1.82 ms | 0 | 4312.4 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.22 ms | 1.41 ms | 1.70 ms | 2.23 ms | 4.57 ms | 0 | 3372.0 |

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
