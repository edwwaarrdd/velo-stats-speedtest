# velo-stats-speedtest

A benchmark harness that compares the six velo-stats backend implementations -
Go, Java, NestJS, Laravel, Symfony and Django - on the four HTTP endpoints they
all expose. Each of them serves the same data: a ride export enriched with weather
from the free Open-Meteo archive and cycling distances from the public OSRM
routing API.

It starts one backend at a time with Docker Compose, waits for its health check,
measures every endpoint, tears the containers down, and moves on to the next.
They cannot run together anyway, since all six publish port 8000, and running
them alone is also the only way to compare them fairly.

## The two profiles

Each backend is measured in two shapes.

**Development** is each repo's own `docker-compose.yml`: the stack a contributor
runs day to day. Source and database are bind-mounted from the host, no CPU
limits apply, and two of the six backends are deliberately running development
servers.

**Production** is the `docker-compose.prod.yml` added alongside it. Nothing is
bind-mounted, the already-seeded SQLite database is baked into the image at build
time, and each backend runs a real application server:

| Backend | Development | Production |
| --- | --- | --- |
| Go | `velo serve`, all host cores | same binary, `GOMAXPROCS=4` |
| Java | embedded server, all host cores | same jar, `-XX:ActiveProcessorCount=4` |
| NestJS | one Node process | four clustered Node processes |
| Laravel | `php artisan serve`, no OPcache | nginx to a four-worker PHP-FPM pool, OPcache with JIT |
| Symfony | `php -S`, no OPcache | nginx to a four-worker PHP-FPM pool, OPcache with JIT |
| Django | gunicorn, one worker, `--reload` | gunicorn, four preloaded workers |

Every production container is capped at four CPUs, and every single-threaded
runtime is given four workers to match, so the languages are compared on an equal
share of the machine rather than on however many cores each happens to grab.

## Reports

A full run writes three files:

| File | Contents |
| --- | --- |
| `RESULTS-DEV.md` | The six backends ranked, development stacks |
| `RESULTS-PROD.md` | The six backends ranked, production stacks |
| `RESULTS-DEV-VS-PROD.md` | Each backend's two profiles side by side |

## Running it

The harness expects the other repos to sit next to this one, which is the default
layout:

```
velo-stats/
  velo-stats-golang/
  velo-stats-laravel/
  velo-stats-nodejs/
  velo-stats-python/
  velo-stats-symfony/
  velo-stats-speedtest/   <- run from here
```

With Go installed:

```bash
go run .
```

Without Go installed, build a native binary in a container once and run that:

```bash
docker run --rm -v "$PWD":/src -w /src -e CGO_ENABLED=0 -e GOOS=darwin -e GOARCH=arm64 golang:1.25-alpine go build -o speedtest .
```

Then `./speedtest`.

A full run builds twelve images and issues thousands of requests per backend per
profile, so it takes a while. To iterate quickly, narrow it:

```bash
go run . -only laravel -profile prod -requests 10 -warmup 2 -skip-build
```

## Flags

| Flag | Default | What it does |
| --- | --- | --- |
| `-requests` | 100 | Requests per endpoint, in each of the two passes. The published results use 1000 |
| `-concurrency` | 10 | Parallel requests during the concurrent pass |
| `-warmup` | 20 | Discarded requests per endpoint before measuring. The published results use 2000, which is what the JVM needs |
| `-only` | all | Run one backend: `golang`, `java`, `nodejs`, `laravel`, `symfony` or `python` |
| `-profile` | `both` | Which stack to measure: `dev`, `prod` or `both` |
| `-skip-build` | false | Reuse the existing image instead of `--build` |
| `-out` | `.` | Directory to write the reports into |
| `-boot-timeout` | 5m | How long to wait for a backend to answer its health check |

## What it measures

For every endpoint on every backend:

1. **Preflight** - one request, recording the status code and body size. This is
   what catches an unseeded database or an unexpected redirect.
2. **Warmup** - discarded requests, so JIT compilation, PHP opcode caching,
   connection setup and the SQLite page cache are not charged to the first
   measurement.
3. **Sequential pass** - requests one at a time. This is per-request latency,
   reported as min, median, p95, p99 and max.
4. **Concurrent pass** - the same number of requests spread over several
   workers. This is throughput, reported as requests per second, and it is where
   a single-process server separates from a concurrent one.

Only the service that exposes the HTTP port is started, plus whatever it depends
on. The three queue workers each repo defines are left down: they do no work
during a benchmark, but they compete for CPU while they boot.

## Source layout

| File | Contents |
| --- | --- |
| `main.go` | Flags and the loop over backends and profiles |
| `backends.go` | The six backend definitions, both profiles, and the endpoint list |
| `docker.go` | Compose wrappers and the health-check poll |
| `bench.go` | The sequential and concurrent passes |
| `stats.go` | Percentiles and throughput |
| `report.go` | Ranking, comparison and rendering |
| `templates.go` | The Markdown templates for the three reports |

Standard library only, no dependencies.

## Results

Production stacks only, from the run recorded in `RESULTS-PROD.md`. **Every
figure is the result of 1000 requests per endpoint, in each of two passes**, on
top of 2000 discarded warmup requests per endpoint. Median latency is the median
of a backend's four per-endpoint medians in the sequential pass; throughput is
the sum of its four endpoints' concurrent requests per second. Fastest first.

| # | Backend | Stack | Median latency | Throughput (req/s) |
| --- | --- | --- | --- | --- |
| 1 | Go | Go 1.25 | 1.14 ms | 8358 |
| 2 | Java | Spring Boot 4.1 / Java 25 | 1.17 ms | 31099 |
| 3 | NestJS | NestJS 11 / Node 22 | 1.56 ms | 19903 |
| 4 | Django | Django 6 / gunicorn | 1.65 ms | 10382 |
| 5 | Symfony | Symfony 8.1 / PHP 8.5 | 3.34 ms | 9315 |
| 6 | Laravel | Laravel 13 / PHP 8.5 | 10.19 ms | 3840 |

Two things that table hides, both worth knowing before quoting it:

**Go and Java are not really separated by latency.** The 0.03 ms between them is
an artefact of how the summary metric picks a median from four values. Endpoint
by endpoint, Java is the faster of the two on three of the four, and Go is ahead
only on `/stations`, by 0.03 ms.

**Throughput and latency disagree about the winner.** Java serves the most
requests per second on three of the four endpoints and roughly four times Go's
total, because a warm JVM with four cores parallelises better than one Go
process under the same cap. Go wins on serving a single request on an idle
server. Which of those matters is a question about the traffic, not about the
language.

The full tables, including percentiles and response sizes, are in
`RESULTS-PROD.md`, and `RESULTS-DEV-VS-PROD.md` sets each backend's development
stack against its production one.

## Caveats

The reports repeat these next to the numbers, because they matter more than the
numbers do:

- Development numbers are uncapped and production numbers are capped at four
  CPUs. A backend that used all twelve host cores in development can look slower
  under the production cap without anything being wrong.
- All six read SQLite. In development that file is bind-mounted, so on macOS
  every read goes through Docker Desktop's filesystem layer; in production it is
  inside the image and that layer is gone. Part of the production gain is this,
  not the application server.
- One run on one machine is not a conclusion. Re-run before trusting a small gap.
- The ranking compares six implementations, not six frameworks. Most of the
  spread is per-row object construction rather than framework overhead, which
  the closing section of each report sets out with measurements.
