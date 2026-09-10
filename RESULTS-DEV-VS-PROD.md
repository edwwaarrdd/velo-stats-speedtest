# Velo-stats backends: development stack versus production stack

Generated 2026-09-10 11:15:41 CEST on darwin/arm64, 12 CPUs with Docker version 29.7.2, build a7dcaa6.

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
| symfony | 4.24 ms | 1.18 ms | 3.60x faster | 1657.5 | 16503.2 | 9.96x faster |
| php | 7.47 ms | 3.93 ms | 1.90x faster | 860.7 | 7857.7 | 9.13x faster |
| python | 1.68 ms | 1.47 ms | 1.14x faster | 4783.3 | 17175.2 | 3.59x faster |
| nodejs | 0.86 ms | 0.80 ms | 1.07x faster | 10046.1 | 24713.7 | 2.46x faster |
| golang | 0.56 ms | 0.52 ms | 1.07x faster | 23680.6 | 27056.4 | 1.14x faster |

## Throughput change per endpoint

Production requests per second divided by development requests per second, for
each endpoint.

| Endpoint | golang | nodejs | php | symfony | python |
| --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 1.24x faster | 1.81x faster | 9.90x faster | 8.09x faster | 2.93x faster |
| `/rides` | 1.08x slower | 2.13x faster | 5.23x faster | 7.03x faster | 5.48x faster |
| `/rides/summary` | 1.08x slower | 3.17x faster | 10.99x faster | 13.91x faster | 4.39x faster |
| `/rides/cost` | 1.01x slower | 2.58x faster | 6.33x faster | 13.43x faster | 4.34x faster |
| `/stations` | 1.02x faster | 3.38x faster | 6.61x faster | 8.58x faster | 4.04x faster |

## Median latency side by side

| Endpoint | golang dev | golang prod | nodejs dev | nodejs prod | php dev | php prod | symfony dev | symfony prod | python dev | python prod |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `/_healthcheck` | 0.23 ms | 0.24 ms | 0.38 ms | 0.33 ms | 2.84 ms | 1.14 ms | 1.53 ms | 0.59 ms | 0.71 ms | 0.57 ms |
| `/rides` | 2.94 ms | 3.07 ms | 4.07 ms | 3.87 ms | 21.13 ms | 15.38 ms | 12.00 ms | 5.75 ms | 8.29 ms | 5.97 ms |
| `/rides/summary` | 0.56 ms | 0.52 ms | 0.43 ms | 0.39 ms | 4.79 ms | 1.57 ms | 4.01 ms | 1.00 ms | 2.07 ms | 1.83 ms |
| `/rides/cost` | 0.43 ms | 0.48 ms | 0.86 ms | 0.80 ms | 15.97 ms | 9.65 ms | 4.24 ms | 1.18 ms | 1.52 ms | 1.14 ms |
| `/stations` | 1.05 ms | 1.09 ms | 1.73 ms | 1.67 ms | 7.47 ms | 3.93 ms | 6.07 ms | 2.52 ms | 1.68 ms | 1.47 ms |

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
