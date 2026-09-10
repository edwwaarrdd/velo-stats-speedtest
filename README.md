# velo-stats-speedtest

A benchmark harness that compares the five velo-stats backend implementations -
Go, NestJS, Laravel, Symfony and Django - on the five HTTP endpoints they all
expose.

It starts one backend at a time with Docker Compose, waits for its health check,
measures every endpoint, tears the containers down, and moves on to the next.
They cannot run together anyway, since all four publish port 8000, and running
them alone is also the only way to compare them fairly.

## The two profiles

Each backend is measured in two shapes.

**Development** is each repo's own `docker-compose.yml`: the stack a contributor
runs day to day. Source and database are bind-mounted from the host, no CPU
limits apply, and two of the four backends are deliberately running development
servers.

**Production** is the `docker-compose.prod.yml` added alongside it. Nothing is
bind-mounted, the already-seeded SQLite database is baked into the image at build
time, and each backend runs a real application server:

| Backend | Development | Production |
| --- | --- | --- |
| Go | `velo serve`, all host cores | same binary, `GOMAXPROCS=4` |
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
| `RESULTS-DEV.md` | The four backends ranked, development stacks |
| `RESULTS-PROD.md` | The four backends ranked, production stacks |
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

A full run builds ten images and issues a thousand requests per backend per
profile, so it takes a while. To iterate quickly, narrow it:

```bash
go run . -only php -profile prod -requests 10 -warmup 2 -skip-build
```

## Flags

| Flag | Default | What it does |
| --- | --- | --- |
| `-requests` | 100 | Requests per endpoint, in each of the two passes |
| `-concurrency` | 10 | Parallel requests during the concurrent pass |
| `-warmup` | 20 | Discarded requests per endpoint before measuring |
| `-only` | all | Run one backend: `golang`, `nodejs`, `php`, `symfony` or `python` |
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
| `backends.go` | The five backend definitions, both profiles, and the endpoint list |
| `docker.go` | Compose wrappers and the health-check poll |
| `bench.go` | The sequential and concurrent passes |
| `stats.go` | Percentiles and throughput |
| `report.go` | Ranking, comparison and rendering |
| `templates.go` | The Markdown templates for the three reports |

Standard library only, no dependencies.

## Caveats

The reports repeat these next to the numbers, because they matter more than the
numbers do:

- Development numbers are uncapped and production numbers are capped at four
  CPUs. A backend that used all twelve host cores in development can look slower
  under the production cap without anything being wrong.
- All four read SQLite. In development that file is bind-mounted, so on macOS
  every read goes through Docker Desktop's filesystem layer; in production it is
  inside the image and that layer is gone. Part of the production gain is this,
  not the application server.
- One run on one machine is not a conclusion. Re-run before trusting a small gap.
