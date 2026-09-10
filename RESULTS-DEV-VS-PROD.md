# Velo-stats backends: development stack versus production stack

Generated 2026-09-10 13:38:18 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| Requests per endpoint per pass | 500 |
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
| symfony | 5.83 ms | 2.56 ms | 2.28x faster | 846.3 | 10523.6 | 12.44x faster |
| laravel | 14.89 ms | 9.84 ms | 1.51x faster | 511.0 | 4199.0 | 8.22x faster |
| python | 1.93 ms | 1.64 ms | 1.17x faster | 2305.3 | 11237.0 | 4.87x faster |
| nodejs | 1.56 ms | 1.67 ms | 1.07x slower | 5880.9 | 16426.9 | 2.79x faster |
| golang | 1.06 ms | 1.09 ms | 1.02x slower | 8559.4 | 8690.9 | 1.02x faster |

## Throughput change per endpoint

Production requests per second divided by development requests per second, for
each endpoint.

| Endpoint | golang | nodejs | laravel | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/rides` | 1.03x faster | 2.84x faster | 5.12x faster | 7.16x faster | 3.95x faster |
| `/rides/summary` | 1.02x faster | 3.01x faster | 10.80x faster | 14.46x faster | 4.96x faster |
| `/rides/cost` | 1.01x faster | 2.05x faster | 5.77x faster | 14.67x faster | 5.23x faster |
| `/stations` | 1.01x faster | 3.29x faster | 6.21x faster | 8.52x faster | 4.61x faster |

## Median latency side by side

| Endpoint | golang dev | golang prod | nodejs dev | nodejs prod | laravel dev | laravel prod | symfony dev | symfony prod | python dev | python prod |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `/rides` | 3.00 ms | 2.95 ms | 3.98 ms | 3.88 ms | 20.46 ms | 15.20 ms | 12.04 ms | 5.52 ms | 6.46 ms | 5.74 ms |
| `/rides/summary` | 0.53 ms | 0.53 ms | 0.40 ms | 0.41 ms | 4.07 ms | 1.60 ms | 3.58 ms | 0.98 ms | 1.93 ms | 1.64 ms |
| `/rides/cost` | 0.46 ms | 0.43 ms | 0.79 ms | 0.77 ms | 14.89 ms | 9.84 ms | 3.79 ms | 1.14 ms | 1.42 ms | 1.07 ms |
| `/stations` | 1.06 ms | 1.09 ms | 1.56 ms | 1.67 ms | 6.83 ms | 4.01 ms | 5.83 ms | 2.56 ms | 1.63 ms | 1.41 ms |

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
