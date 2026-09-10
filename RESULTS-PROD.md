# Velo-stats backend speed test - production stack

Generated 2026-09-10 15:47:29 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

Each backend was started on its own, warmed up, measured, then torn down before
the next one started. No two backends ran at the same time.

**Every figure below is based on 1000 requests per endpoint,
in each of two passes.** A further 2000 requests per endpoint
are issued first and discarded, which is set by the slowest runtime to warm
rather than the fastest: the JVM needs thousands of calls to finish compiling,
where the others are at full speed within tens.

This is the **production** profile: each backend's `docker-compose.prod.yml`,
running a real application server with no bind mounts and the already-seeded
SQLite database baked into the image. Every backend is capped at the same four
CPUs, so the runtimes are compared on an equal share of the machine.

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
| 1 | golang | Go 1.25 | net/http, baked-in SQLite, GOMAXPROCS=4 | 1.14 ms | 8358.2 | 0 |
| 2 | java | Spring Boot 4.1 / Java 25 | embedded server, baked-in SQLite, ActiveProcessorCount=4 | 1.17 ms | 31098.5 | 0 |
| 3 | nodejs | NestJS 11 / Node 22 | 4 clustered Node processes, baked-in SQLite | 1.56 ms | 19902.7 | 0 |
| 4 | python | Django 6 / gunicorn | gunicorn, 4 preloaded workers, baked-in SQLite | 1.65 ms | 10381.9 | 0 |
| 5 | symfony | Symfony 8.1 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 3.34 ms | 9315.4 | 0 |
| 6 | laravel | Laravel 13 / PHP 8.5 | nginx to a 4-worker PHP-FPM pool, OPcache with JIT | 10.19 ms | 3840.2 | 0 |

## Median latency per endpoint (sequential)

| Endpoint | golang | java | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- | --- |
| `/rides` | 3.09 ms | 2.53 ms | 3.93 ms | 15.21 ms | 6.27 ms | 5.99 ms |
| `/rides/summary` | 0.52 ms | 0.34 ms | 0.35 ms | 1.64 ms | 1.00 ms | 1.65 ms |
| `/rides/cost` | 0.44 ms | 0.32 ms | 0.72 ms | 10.19 ms | 1.25 ms | 1.07 ms |
| `/stations` | 1.14 ms | 1.17 ms | 1.56 ms | 4.54 ms | 3.34 ms | 1.38 ms |

## Throughput per endpoint (concurrent, req/s)

| Endpoint | golang | java | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- | --- |
| `/rides` | 503.5 | 876.4 | 793.9 | 248.8 | 582.1 | 600.6 |
| `/rides/summary` | 2882.9 | 11443.3 | 11699.8 | 2404.2 | 4103.6 | 2769.7 |
| `/rides/cost` | 3725.6 | 15378.5 | 4813.5 | 402.2 | 3390.1 | 4371.6 |
| `/stations` | 1246.1 | 3400.3 | 2595.5 | 784.9 | 1239.5 | 2640.0 |

## Detail per backend

### golang - Go 1.25

net/http, baked-in SQLite, GOMAXPROCS=4 (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.89 ms | 3.09 ms | 3.43 ms | 3.69 ms | 6.73 ms | 0 | 503.5 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.44 ms | 0.52 ms | 0.58 ms | 0.64 ms | 2.34 ms | 0 | 2882.9 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.38 ms | 0.44 ms | 0.76 ms | 1.62 ms | 2.31 ms | 0 | 3725.6 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.01 ms | 1.14 ms | 1.42 ms | 1.52 ms | 2.08 ms | 0 | 1246.1 |

### java - Spring Boot 4.1 / Java 25

embedded server, baked-in SQLite, ActiveProcessorCount=4 (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 2.23 ms | 2.53 ms | 3.21 ms | 3.74 ms | 4.45 ms | 0 | 876.4 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.27 ms | 0.34 ms | 0.44 ms | 0.91 ms | 1.51 ms | 0 | 11443.3 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 0.25 ms | 0.32 ms | 0.39 ms | 0.86 ms | 1.66 ms | 0 | 15378.5 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.05 ms | 1.17 ms | 1.33 ms | 2.95 ms | 6.51 ms | 0 | 3400.3 |

### nodejs - NestJS 11 / Node 22

4 clustered Node processes, baked-in SQLite (`docker-compose.prod.yml`, service `app`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 123.0 KB | 3.73 ms | 3.93 ms | 4.41 ms | 4.69 ms | 5.69 ms | 0 | 793.9 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.29 ms | 0.35 ms | 0.44 ms | 0.56 ms | 2.06 ms | 0 | 11699.8 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 335 B | 0.63 ms | 0.72 ms | 0.83 ms | 1.00 ms | 2.17 ms | 0 | 4813.5 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 1.39 ms | 1.56 ms | 1.66 ms | 1.84 ms | 2.96 ms | 0 | 2595.5 |

### laravel - Laravel 13 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 14.86 ms | 15.21 ms | 16.21 ms | 18.37 ms | 30.91 ms | 0 | 248.8 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 1.48 ms | 1.64 ms | 2.03 ms | 2.35 ms | 2.66 ms | 0 | 2404.2 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 9.41 ms | 10.19 ms | 11.50 ms | 13.29 ms | 28.03 ms | 0 | 402.2 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 3.82 ms | 4.54 ms | 5.51 ms | 8.24 ms | 12.62 ms | 0 | 784.9 |

### symfony - Symfony 8.1 / PHP 8.5

nginx to a 4-worker PHP-FPM pool, OPcache with JIT (`docker-compose.prod.yml`, service `web`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides` | 200 | 125.2 KB | 5.35 ms | 6.27 ms | 7.69 ms | 9.07 ms | 13.36 ms | 0 | 582.1 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 186 B | 0.84 ms | 1.00 ms | 1.51 ms | 2.33 ms | 6.85 ms | 0 | 4103.6 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 341 B | 1.05 ms | 1.25 ms | 1.53 ms | 4.27 ms | 4.74 ms | 0 | 3390.1 |
| `http://127.0.0.1:8000/stations` | 200 | 24.2 KB | 2.94 ms | 3.34 ms | 3.99 ms | 5.65 ms | 8.54 ms | 0 | 1239.5 |

### python - Django 6 / gunicorn

gunicorn, 4 preloaded workers, baked-in SQLite (`docker-compose.prod.yml`, service `api`)

| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `http://127.0.0.1:8000/rides/` | 200 | 134.7 KB | 5.35 ms | 5.99 ms | 6.81 ms | 7.42 ms | 16.70 ms | 0 | 600.6 |
| `http://127.0.0.1:8000/rides/summary` | 200 | 199 B | 1.38 ms | 1.65 ms | 1.97 ms | 2.73 ms | 4.89 ms | 0 | 2769.7 |
| `http://127.0.0.1:8000/rides/cost` | 200 | 362 B | 0.89 ms | 1.07 ms | 1.33 ms | 1.53 ms | 3.29 ms | 0 | 4371.6 |
| `http://127.0.0.1:8000/stations/` | 200 | 26.8 KB | 1.20 ms | 1.38 ms | 1.66 ms | 2.08 ms | 7.03 ms | 0 | 2640.0 |

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
  workers, four clustered Node processes, `GOMAXPROCS=4` for Go and
  `-XX:ActiveProcessorCount=4` for the JVM. Without that cap Go and the JVM
  would simply take all twelve cores of the host.
- **Nothing is bind-mounted.** The seeded SQLite database is copied into the
  image at build time, so no request touches the host filesystem.
- **This is still a single container per backend on a laptop.** It is not a
  tuned deployment, and it says nothing about how these stacks behave behind a
  load balancer, with a networked database, or under sustained traffic.

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
