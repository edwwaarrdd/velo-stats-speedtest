# Velo-stats backend speed test - production stack

Generated 2026-09-10 10:46:16 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| 1 | golang | Go 1.25 | net/http, baked-in SQLite, GOMAXPROCS=4 | 0.52 ms | 24185.5 | 0 |
| 2 | symfony | Symfony 8.1 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 1.16 ms | 18440.8 | 0 |
| 3 | python | Django 6 / gunicorn | gunicorn, 4 preloaded workers, baked-in SQLite | 1.36 ms | 20979.4 | 0 |
| 4 | nodejs | NestJS 11 / Node 22 | 4 clustered Node processes, baked-in SQLite | 1.67 ms | 15746.4 | 0 |
| 5 | php | Laravel 13 / PHP 8.4 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 4.41 ms | 6407.7 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | php | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.32 ms | 0.38 ms | 1.21 ms | 0.68 ms | 0.66 ms |
| `/rides` | 3.07 ms | 4.63 ms | 17.38 ms | 3.06 ms | 6.06 ms |
| `/rides/summary` | 0.52 ms | 1.35 ms | 1.75 ms | 1.11 ms | 1.69 ms |
| `/rides/cost` | 0.47 ms | 1.67 ms | 10.99 ms | 1.16 ms | 1.11 ms |
| `/stations` | 1.11 ms | 1.83 ms | 4.41 ms | 1.29 ms | 1.36 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | php | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 15702.6 | 5705.1 | 2753.9 | 6673.1 | 8254.6 |
| `/rides` | 499.1 | 363.0 | 221.2 | 1345.5 | 642.9 |
| `/rides/summary` | 3055.0 | 3540.6 | 2334.1 | 3661.4 | 2878.9 |
| `/rides/cost` | 3629.5 | 3945.1 | 366.1 | 3615.2 | 5164.7 |
| `/stations` | 1299.3 | 2192.6 | 732.5 | 3145.5 | 4038.3 |

## Detail per backend

### golang - Go 1.25

net/http, baked-in SQLite, GOMAXPROCS=4 (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.18 ms | 0.32 ms | 0.52 ms | 1.41 ms | 1.41 ms | 0 | 15702.6 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.85 ms | 3.07 ms | 3.78 ms | 4.63 ms | 4.63 ms | 0 | 499.1 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.46 ms | 0.52 ms | 0.61 ms | 0.65 ms | 0.65 ms | 0 | 3055.0 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.41 ms | 0.47 ms | 0.57 ms | 0.85 ms | 0.85 ms | 0 | 3629.5 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.97 ms | 1.11 ms | 1.45 ms | 1.72 ms | 1.72 ms | 0 | 1299.3 |

### nodejs - NestJS 11 / Node 22

4 clustered Node processes, baked-in SQLite (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.30 ms | 0.38 ms | 0.50 ms | 0.70 ms | 0.70 ms | 0 | 5705.1 |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.91 ms | 4.63 ms | 5.25 ms | 5.70 ms | 5.70 ms | 0 | 363.0 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.42 ms | 1.35 ms | 7.67 ms | 34.95 ms | 34.95 ms | 0 | 3540.6 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.86 ms | 1.67 ms | 3.80 ms | 4.24 ms | 4.24 ms | 0 | 3945.1 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.70 ms | 1.83 ms | 2.06 ms | 2.17 ms | 2.17 ms | 0 | 2192.6 |

### php - Laravel 13 / PHP 8.4

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 1.06 ms | 1.21 ms | 1.50 ms | 1.92 ms | 1.92 ms | 0 | 2753.9 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 15.73 ms | 17.38 ms | 20.85 ms | 50.32 ms | 50.32 ms | 0 | 221.2 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 1.55 ms | 1.75 ms | 1.99 ms | 2.03 ms | 2.03 ms | 0 | 2334.1 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 10.23 ms | 10.99 ms | 11.67 ms | 11.81 ms | 11.81 ms | 0 | 366.1 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 4.08 ms | 4.41 ms | 4.93 ms | 5.07 ms | 5.07 ms | 0 | 732.5 |

### symfony - Symfony 8.1 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.58 ms | 0.68 ms | 0.85 ms | 1.05 ms | 1.05 ms | 0 | 6673.1 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.77 ms | 3.06 ms | 3.60 ms | 3.97 ms | 3.97 ms | 0 | 1345.5 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 1.00 ms | 1.11 ms | 1.29 ms | 1.49 ms | 1.49 ms | 0 | 3661.4 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.99 ms | 1.16 ms | 1.35 ms | 1.65 ms | 1.65 ms | 0 | 3615.2 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.17 ms | 1.29 ms | 1.48 ms | 1.50 ms | 1.50 ms | 0 | 3145.5 |

### python - Django 6 / gunicorn

gunicorn, 4 preloaded workers, baked-in SQLite (`docker-compose.prod.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 17 B | 0.54 ms | 0.66 ms | 1.48 ms | 2.36 ms | 2.36 ms | 0 | 8254.6 |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.58 ms | 6.06 ms | 6.83 ms | 10.97 ms | 10.97 ms | 0 | 642.9 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.43 ms | 1.69 ms | 1.99 ms | 2.58 ms | 2.58 ms | 0 | 2878.9 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 0.94 ms | 1.11 ms | 1.31 ms | 1.40 ms | 1.40 ms | 0 | 5164.7 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.21 ms | 1.36 ms | 1.57 ms | 1.62 ms | 1.62 ms | 0 | 4038.3 |

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
