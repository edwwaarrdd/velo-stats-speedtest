# Velo-stats backends: development stack versus production stack

Generated 2026-09-10 15:47:29 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| Requests per endpoint per pass | 1000 |
| Concurrency in the concurrent pass | 10 |
| Discarded warmup requests per endpoint | 2000 |

## What changed per backend

| Backend | Development | Production |
| --- | --- | --- |
| symfony | php -S, one request at a time, no OPcache | nginx to a 4-worker PHP-FPM pool, OPcache with JIT |
| laravel | php artisan serve, one request at a time, no OPcache | nginx to a 4-worker PHP-FPM pool, OPcache with JIT |
| python | gunicorn, 1 worker, --reload, bind-mounted source | gunicorn, 4 preloaded workers, baked-in SQLite |
| nodejs | single Node process, bind-mounted source and SQLite | 4 clustered Node processes, baked-in SQLite |
| golang | net/http, bind-mounted SQLite, all cores | net/http, baked-in SQLite, GOMAXPROCS=4 |
| java | embedded server, bind-mounted SQLite, all cores | embedded server, baked-in SQLite, ActiveProcessorCount=4 |

## The gain, biggest first

| Backend | Dev latency | Prod latency | Latency change | Dev req/s | Prod req/s | Throughput change |
| --- | --- | --- | --- | --- | --- | --- |
| symfony | 5.64 ms | 3.34 ms | 1.69x faster | 845.3 | 9315.4 | 11.02x faster |
| laravel | 15.06 ms | 10.19 ms | 1.48x faster | 510.9 | 3840.2 | 7.52x faster |
| python | 2.02 ms | 1.65 ms | 1.23x faster | 2212.1 | 10381.9 | 4.69x faster |
| nodejs | 1.58 ms | 1.56 ms | 1.01x faster | 6181.6 | 19902.7 | 3.22x faster |
| golang | 1.12 ms | 1.14 ms | 1.02x slower | 8287.4 | 8358.2 | 1.01x faster |
| java | 1.15 ms | 1.17 ms | 1.01x slower | 32379.6 | 31098.5 | 1.04x slower |

## Throughput change per endpoint

Production requests per second divided by development requests per second, for
each endpoint.

| Endpoint | golang | java | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- | --- |
| `/rides` | 1.00x slower | 1.47x slower | 3.05x faster | 5.24x faster | 7.16x faster | 3.64x faster |
| `/rides/summary` | 1.07x faster | 1.13x slower | 3.19x faster | 9.29x faster | 13.86x faster | 4.78x faster |
| `/rides/cost` | 1.01x slower | 1.09x faster | 3.07x faster | 6.02x faster | 11.97x faster | 5.95x faster |
| `/stations` | 1.05x slower | 1.17x slower | 3.79x faster | 5.70x faster | 6.71x faster | 3.60x faster |

## Median latency side by side

| Endpoint | golang dev | golang prod | java dev | java prod | nodejs dev | nodejs prod | laravel dev | laravel prod | symfony dev | symfony prod | python dev | python prod |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `/rides` | 2.93 ms | 3.09 ms | 2.48 ms | 2.53 ms | 3.91 ms | 3.93 ms | 20.84 ms | 15.21 ms | 11.53 ms | 6.27 ms | 6.35 ms | 5.99 ms |
| `/rides/summary` | 0.52 ms | 0.52 ms | 0.33 ms | 0.34 ms | 0.35 ms | 0.35 ms | 4.18 ms | 1.64 ms | 3.95 ms | 1.00 ms | 2.02 ms | 1.65 ms |
| `/rides/cost` | 0.42 ms | 0.44 ms | 0.31 ms | 0.32 ms | 0.74 ms | 0.72 ms | 15.06 ms | 10.19 ms | 3.88 ms | 1.25 ms | 1.48 ms | 1.07 ms |
| `/stations` | 1.12 ms | 1.14 ms | 1.15 ms | 1.17 ms | 1.58 ms | 1.56 ms | 8.28 ms | 4.54 ms | 5.64 ms | 3.34 ms | 1.72 ms | 1.38 ms |

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
