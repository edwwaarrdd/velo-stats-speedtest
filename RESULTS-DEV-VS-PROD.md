# Velo-stats backends: development stack versus production stack

Generated 2026-09-10 10:46:16 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| php | php artisan serve, one request at a time, no OPcache | nginx to a 4-worker PHP-FPM pool, OPcache with JIT |
| python | gunicorn, 1 worker, --reload, bind-mounted source | gunicorn, 4 preloaded workers, baked-in SQLite |
| nodejs | single Node process, bind-mounted source and SQLite | 4 clustered Node processes, baked-in SQLite |
| golang | net/http, bind-mounted SQLite, all cores | net/http, baked-in SQLite, GOMAXPROCS=4 |

## The gain, biggest first

| Backend | Dev latency | Prod latency | Latency change | Dev req/s | Prod req/s | Throughput change |
| --- | --- | --- | --- | --- | --- | --- |
| symfony | 3.63 ms | 1.16 ms | 3.12x faster | 1891.1 | 18440.8 | 9.75x faster |
| php | 7.18 ms | 4.41 ms | 1.63x faster | 855.2 | 6407.7 | 7.49x faster |
| python | 1.61 ms | 1.36 ms | 1.19x faster | 5469.4 | 20979.4 | 3.84x faster |
| nodejs | 0.83 ms | 1.67 ms | 2.02x slower | 9767.5 | 15746.4 | 1.61x faster |
| golang | 0.54 ms | 0.52 ms | 1.04x faster | 22859.8 | 24185.5 | 1.06x faster |

## Throughput change per endpoint

Production requests per second divided by development requests per second, for
each endpoint.

| Endpoint | golang | nodejs | php | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 1.10x faster | 1.29x faster | 8.49x faster | 8.05x faster | 2.69x faster |
| `/rides` | 1.02x faster | 1.44x faster | 4.77x faster | 8.12x faster | 3.84x faster |
| `/rides/summary` | 1.01x slower | 1.17x faster | 8.72x faster | 12.59x faster | 4.71x faster |
| `/rides/cost` | 1.02x slower | 2.79x faster | 5.55x faster | 11.56x faster | 5.87x faster |
| `/stations` | 1.03x faster | 3.37x faster | 4.86x faster | 10.74x faster | 5.45x faster |

## Median latency side by side

| Endpoint | golang dev | golang prod | nodejs dev | nodejs prod | php dev | php prod | symfony dev | symfony prod | python dev | python prod |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.29 ms | 0.32 ms | 0.42 ms | 0.38 ms | 2.74 ms | 1.21 ms | 1.68 ms | 0.68 ms | 0.58 ms | 0.66 ms |
| `/rides` | 3.17 ms | 3.07 ms | 4.14 ms | 4.63 ms | 22.94 ms | 17.38 ms | 7.53 ms | 3.06 ms | 6.20 ms | 6.06 ms |
| `/rides/summary` | 0.54 ms | 0.52 ms | 0.42 ms | 1.35 ms | 4.30 ms | 1.75 ms | 4.13 ms | 1.11 ms | 1.95 ms | 1.69 ms |
| `/rides/cost` | 0.45 ms | 0.47 ms | 0.83 ms | 1.67 ms | 16.39 ms | 10.99 ms | 3.49 ms | 1.16 ms | 1.44 ms | 1.11 ms |
| `/stations` | 1.09 ms | 1.11 ms | 1.65 ms | 1.83 ms | 7.18 ms | 4.41 ms | 3.63 ms | 1.29 ms | 1.61 ms | 1.36 ms |

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
