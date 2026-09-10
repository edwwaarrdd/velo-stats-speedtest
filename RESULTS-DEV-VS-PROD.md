# Velo-stats backends: development stack versus production stack

Generated 2026-09-10 09:13:06 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

Both stacks were measured in the same run, minutes apart, on the same machine.
The development profile is each repo's own `docker-compose.yml`. The production
profile is the `docker-compose.prod.yml` alongside it: a real application server,
no bind mounts, the already-seeded SQLite database baked into the image, and a
four-CPU cap applied equally to every backend.

The point of this comparison is to separate two things that the development
numbers alone conflate: how fast a language and framework are, and how much of a
backend's measured slowness was really its development server.

| Setting | Value |
| --- | --- |
| Requests per endpoint per pass | 100 |
| Concurrency in the concurrent pass | 10 |
| Discarded warmup requests per endpoint | 20 |

## What changed per backend

| Backend | Development | Production |
| --- | --- | --- |
| php | php artisan serve, one request at a time, no OPcache | nginx to a 4-worker PHP-FPM pool, OPcache with JIT |
| python | gunicorn, 1 worker, --reload, bind-mounted source | gunicorn, 4 preloaded workers, baked-in SQLite |
| nodejs | single Node process, bind-mounted source and SQLite | 4 clustered Node processes, baked-in SQLite |
| golang | net/http, bind-mounted SQLite, all cores | net/http, baked-in SQLite, GOMAXPROCS=4 |

## The gain, biggest first

| Backend | Dev latency | Prod latency | Latency change | Dev req/s | Prod req/s | Throughput change |
| --- | --- | --- | --- | --- | --- | --- |
| php | 7.01 ms | 3.98 ms | 1.76x faster | 882.0 | 7344.7 | 8.33x faster |
| python | 1.58 ms | 1.31 ms | 1.20x faster | 5443.9 | 18177.5 | 3.34x faster |
| nodejs | 0.80 ms | 0.76 ms | 1.05x faster | 10222.3 | 28487.5 | 2.79x faster |
| golang | 0.54 ms | 0.50 ms | 1.07x faster | 24901.9 | 26338.4 | 1.06x faster |

## Throughput change per endpoint

Production requests per second divided by development requests per second, for
each endpoint.

| Endpoint | golang | nodejs | php | python |
| --- | --- | --- | --- | --- |
| `/_healthcheck` | 1.03x faster | 2.10x faster | 9.94x faster | 1.64x faster |
| `/rides` | 1.09x faster | 2.18x faster | 4.49x faster | 4.27x faster |
| `/rides/summary` | 1.16x faster | 3.50x faster | 8.44x faster | 5.11x faster |
| `/rides/cost` | 1.06x faster | 3.10x faster | 5.71x faster | 6.70x faster |
| `/stations` | 1.11x faster | 3.81x faster | 6.71x faster | 5.39x faster |

## Median latency side by side

| Endpoint | golang dev | golang prod | nodejs dev | nodejs prod | php dev | php prod | python dev | python prod |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.24 ms | 0.24 ms | 0.35 ms | 0.37 ms | 3.02 ms | 1.15 ms | 0.68 ms | 0.56 ms |
| `/rides` | 3.02 ms | 2.81 ms | 3.85 ms | 3.86 ms | 19.74 ms | 16.12 ms | 6.11 ms | 5.54 ms |
| `/rides/summary` | 0.54 ms | 0.50 ms | 0.42 ms | 0.42 ms | 3.94 ms | 1.73 ms | 1.99 ms | 1.56 ms |
| `/rides/cost` | 0.47 ms | 0.42 ms | 0.80 ms | 0.76 ms | 14.85 ms | 10.15 ms | 1.40 ms | 1.06 ms |
| `/stations` | 1.07 ms | 0.97 ms | 1.65 ms | 1.59 ms | 7.01 ms | 3.98 ms | 1.58 ms | 1.31 ms |

## Reading the gap

- **A large gain means the development stack was the bottleneck, not the
  language.** The clearest case is PHP: `php artisan serve` is a single-process
  development server with no OPcache, so almost any production configuration
  beats it by a wide margin. Judging PHP on its development number is judging
  that server.
- **A small gain means the development stack was already close to production.**
  A backend whose development server already handles requests concurrently and
  already caches its compiled code has little left to win.
- **Latency and throughput move for different reasons.** Caching compiled code
  and dropping bind mounts cuts per-request latency. Adding worker processes
  raises throughput without helping a single idle request at all.
- **The four-CPU cap only applies to the production profile.** Development
  numbers are uncapped, so a backend that used all twelve cores in development
  can legitimately look slower under the production cap. Check the per-endpoint
  table before calling that a regression.
- **Full detail for each profile** is in `RESULTS-DEV.md` and
  `RESULTS-PROD.md`, including per-endpoint percentiles and response sizes.
