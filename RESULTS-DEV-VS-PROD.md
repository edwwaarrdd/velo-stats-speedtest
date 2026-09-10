# Velo-stats backends: development stack versus production stack

Generated 2026-09-10 13:26:04 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| symfony | php -S, one request at a time, no OPcache | nginx to a 4-worker PHP-FPM pool, OPcache with JIT |
| laravel | php artisan serve, one request at a time, no OPcache | nginx to a 4-worker PHP-FPM pool, OPcache with JIT |
| python | gunicorn, 1 worker, --reload, bind-mounted source | gunicorn, 4 preloaded workers, baked-in SQLite |
| nodejs | single Node process, bind-mounted source and SQLite | 4 clustered Node processes, baked-in SQLite |
| golang | net/http, bind-mounted SQLite, all cores | net/http, baked-in SQLite, GOMAXPROCS=4 |

## The gain, biggest first

| Backend | Dev latency | Prod latency | Latency change | Dev req/s | Prod req/s | Throughput change |
| --- | --- | --- | --- | --- | --- | --- |
| symfony | 3.79 ms | 1.12 ms | 3.38x faster | 1744.3 | 16780.6 | 9.62x faster |
| laravel | 6.56 ms | 4.44 ms | 1.48x faster | 989.4 | 7018.5 | 7.09x faster |
| python | 1.77 ms | 1.39 ms | 1.27x faster | 5666.0 | 19276.0 | 3.40x faster |
| nodejs | 0.81 ms | 0.84 ms | 1.04x slower | 11216.3 | 20964.2 | 1.87x faster |
| golang | 0.49 ms | 0.56 ms | 1.14x slower | 25712.4 | 23656.5 | 1.09x slower |

## Throughput change per endpoint

Production requests per second divided by development requests per second, for
each endpoint.

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 1.13x slower | 1.14x faster | 8.03x faster | 7.82x faster | 2.19x faster |
| `/rides` | 1.07x slower | 1.63x faster | 4.64x faster | 7.02x faster | 4.32x faster |
| `/rides/summary` | 1.01x faster | 2.36x faster | 8.06x faster | 13.26x faster | 5.21x faster |
| `/rides/cost` | 1.03x slower | 2.90x faster | 5.01x faster | 12.79x faster | 5.27x faster |
| `/stations` | 1.02x faster | 3.27x faster | 4.52x faster | 8.46x faster | 5.42x faster |

## Median latency side by side

| Endpoint | golang dev | golang prod | nodejs dev | nodejs prod | laravel dev | laravel prod | symfony dev | symfony prod | python dev | python prod |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.22 ms | 0.28 ms | 0.35 ms | 0.35 ms | 2.67 ms | 1.28 ms | 1.45 ms | 0.62 ms | 0.63 ms | 0.61 ms |
| `/rides` | 2.86 ms | 3.27 ms | 3.85 ms | 4.42 ms | 20.72 ms | 16.12 ms | 11.09 ms | 5.78 ms | 6.78 ms | 5.87 ms |
| `/rides/summary` | 0.49 ms | 0.56 ms | 0.41 ms | 0.43 ms | 3.96 ms | 1.78 ms | 3.61 ms | 0.97 ms | 2.01 ms | 1.70 ms |
| `/rides/cost` | 0.43 ms | 0.46 ms | 0.81 ms | 0.84 ms | 14.26 ms | 10.69 ms | 3.79 ms | 1.12 ms | 1.54 ms | 1.14 ms |
| `/stations` | 1.08 ms | 1.11 ms | 1.61 ms | 1.68 ms | 6.56 ms | 4.44 ms | 6.06 ms | 2.52 ms | 1.77 ms | 1.39 ms |

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
