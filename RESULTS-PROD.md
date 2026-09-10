# Velo-stats backend speed test - production stack

Generated 2026-09-10 09:13:06 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| 1 | golang | Go 1.25 | net/http, baked-in SQLite, GOMAXPROCS=4 | 0.50 ms | 26338.4 | 0 |
| 2 | nodejs | NestJS 11 / Node 22 | 4 clustered Node processes, baked-in SQLite | 0.76 ms | 28487.5 | 0 |
| 3 | python | Django 6 / gunicorn | gunicorn, 4 preloaded workers, baked-in SQLite | 1.31 ms | 18177.5 | 0 |
| 4 | php | Laravel 13 / PHP 8.4 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 3.98 ms | 7344.7 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | nodejs | php | python |
| --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.24 ms | 0.37 ms | 1.15 ms | 0.56 ms |
| `/rides` | 2.81 ms | 3.86 ms | 16.12 ms | 5.54 ms |
| `/rides/summary` | 0.50 ms | 0.42 ms | 1.73 ms | 1.56 ms |
| `/rides/cost` | 0.42 ms | 0.76 ms | 10.15 ms | 1.06 ms |
| `/stations` | 0.97 ms | 1.59 ms | 3.98 ms | 1.31 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | nodejs | php | python |
| --- | --- | --- | --- | --- |
| `/_healthcheck` | 17321.3 | 9821.7 | 3522.9 | 5188.6 |
| `/rides` | 534.2 | 578.8 | 234.0 | 719.6 |
| `/rides/summary` | 3322.1 | 11320.3 | 2264.9 | 3064.8 |
| `/rides/cost` | 3761.3 | 4428.2 | 387.0 | 5351.5 |
| `/stations` | 1399.5 | 2338.6 | 936.0 | 3853.0 |

## Detail per backend

### golang - Go 1.25

net/http, baked-in SQLite, GOMAXPROCS=4 (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.17 ms | 0.24 ms | 0.32 ms | 0.35 ms | 0.35 ms | 0 | 17321.3 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.62 ms | 2.81 ms | 3.16 ms | 3.26 ms | 3.26 ms | 0 | 534.2 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.45 ms | 0.50 ms | 0.56 ms | 0.68 ms | 0.68 ms | 0 | 3322.1 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.39 ms | 0.42 ms | 0.48 ms | 0.71 ms | 0.71 ms | 0 | 3761.3 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 0.88 ms | 0.97 ms | 1.23 ms | 1.32 ms | 1.32 ms | 0 | 1399.5 |

### nodejs - NestJS 11 / Node 22

4 clustered Node processes, baked-in SQLite (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 0.29 ms | 0.37 ms | 0.46 ms | 1.92 ms | 1.92 ms | 0 | 9821.7 |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.68 ms | 3.86 ms | 4.20 ms | 4.74 ms | 4.74 ms | 0 | 578.8 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.39 ms | 0.42 ms | 0.47 ms | 0.53 ms | 0.53 ms | 0 | 11320.3 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.72 ms | 0.76 ms | 1.01 ms | 1.41 ms | 1.41 ms | 0 | 4428.2 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.46 ms | 1.59 ms | 1.81 ms | 3.37 ms | 3.37 ms | 0 | 2338.6 |

### php - Laravel 13 / PHP 8.4

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 16 B | 1.03 ms | 1.15 ms | 1.45 ms | 1.86 ms | 1.86 ms | 0 | 3522.9 |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 15.35 ms | 16.12 ms | 17.18 ms | 17.81 ms | 17.81 ms | 0 | 234.0 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 1.53 ms | 1.73 ms | 2.01 ms | 2.34 ms | 2.34 ms | 0 | 2264.9 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 9.30 ms | 10.15 ms | 11.65 ms | 13.19 ms | 13.19 ms | 0 | 387.0 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 3.83 ms | 3.98 ms | 4.27 ms | 5.88 ms | 5.88 ms | 0 | 936.0 |

### python - Django 6 / gunicorn

gunicorn, 4 preloaded workers, baked-in SQLite (`docker-compose.prod.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/_healthcheck` | 200 | 17 B | 0.47 ms | 0.56 ms | 0.70 ms | 0.81 ms | 0.81 ms | 0 | 5188.6 |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.20 ms | 5.54 ms | 5.98 ms | 6.97 ms | 6.97 ms | 0 | 719.6 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.42 ms | 1.56 ms | 1.78 ms | 2.36 ms | 2.36 ms | 0 | 3064.8 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 0.94 ms | 1.06 ms | 1.26 ms | 1.63 ms | 1.63 ms | 0 | 5351.5 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.19 ms | 1.31 ms | 1.52 ms | 1.62 ms | 1.62 ms | 0 | 3853.0 |

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
