package main

// The Markdown templates for the three reports. Keeping them here rather than
// scattered through the rendering code means the shape of a report is visible
// in one place.

const reportTemplate = `# Velo-stats backend speed test - {{ .ProfileName }} stack

Generated {{ .GeneratedAt }} on {{ .Host }} with {{ .Docker }}.

Each backend was started on its own, warmed up, measured, then torn down before
the next one started. No two backends ran at the same time.
{{ if eq .Profile "prod" }}
This is the **production** profile: each backend's ` + "`docker-compose.prod.yml`" + `,
running a real application server with no bind mounts and the already-seeded
SQLite database baked into the image. Every backend is capped at the same four
CPUs, so the runtimes are compared on an equal share of the machine.
{{ else }}
This is the **development** profile: each backend's own ` + "`docker-compose.yml`" + `,
the stack a contributor runs day to day. Source and database are bind-mounted
from the host and no CPU limits are applied.
{{ end }}
| Setting | Value |
| --- | --- |
| Requests per endpoint per pass | {{ .Config.Requests }} |
| Concurrency in the concurrent pass | {{ .Config.Concurrency }} |
| Discarded warmup requests per endpoint | {{ .Config.Warmup }} |
| Endpoints measured | {{ len .Endpoints }} |

## Ranking

Median latency is the median of each backend's five per-endpoint medians in the
sequential pass. Throughput is the sum of the five endpoints' concurrent
requests per second. Fastest first.

| # | Backend | Stack | Setup | Median latency | Throughput (req/s) | Errors |
| --- | --- | --- | --- | --- | --- | --- |
{{ range $i, $r := .Ranking }}| {{ inc $i }} | {{ $r.Name }} | {{ $r.Language }} | {{ $r.Description }} | {{ ms $r.MedianLatency }} | {{ rps $r.TotalRPS }} | {{ $r.Errors }} |
{{ end }}
## Median latency per endpoint (sequential)

| Endpoint |{{ range .Results }} {{ .Backend.Name }} |{{ end }}
| --- |{{ range .Results }} --- |{{ end }}
{{ $results := .Results }}{{ range $endpoint := .Endpoints }}| ` + "`{{ $endpoint }}`" + ` |{{ range $r := $results }} {{ cell $r $endpoint }} |{{ end }}
{{ end }}
## Throughput per endpoint (concurrent, req/s)

| Endpoint |{{ range .Results }} {{ .Backend.Name }} |{{ end }}
| --- |{{ range .Results }} --- |{{ end }}
{{ range $endpoint := .Endpoints }}| ` + "`{{ $endpoint }}`" + ` |{{ range $r := $results }} {{ rpsCell $r $endpoint }} |{{ end }}
{{ end }}
## Detail per backend
{{ range .Results }}
### {{ .Backend.Name }} - {{ .Backend.Language }}

{{ .Deployment.Description }} (` + "`{{ .Deployment.ComposeFile }}`" + `, service ` + "`{{ .Deployment.Service }}`" + `)
{{ if .Failed }}
Not measured: {{ .Note }}
{{ else }}
| Endpoint | Status | Body | Min | Median | p95 | p99 | Max | Errors | Concurrent req/s |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
{{ range .Endpoints }}| ` + "`{{ .URL }}`" + ` | {{ .Status }} | {{ bytes .BodyBytes }} | {{ ms .Sequential.Min }} | {{ ms .Sequential.Median }} | {{ ms .Sequential.P95 }} | {{ ms .Sequential.P99 }} | {{ ms .Sequential.Max }} | {{ .Sequential.Errors }} | {{ rps .Concurrent.RPS }} |
{{ end }}{{ end }}{{ end }}
## How to read this, and what it does not say

- **Compare the body sizes in the detail tables.** If one backend's ` + "`/rides`" + `
  response is far smaller than the others, its database was never seeded and its
  numbers are not comparable.
- **Sequential latency and concurrent throughput answer different questions.**
  The first is how fast one request is served on an idle server; the second is
  how much work the server gets through under load. A runtime can win one and
  lose the other.
- **This is one run on one machine.** Latency at this scale is sensitive to
  background load. Re-run before drawing a conclusion from a small gap.
{{ if eq .Profile "dev" }}- **These are development stacks, and two of them are deliberately slow.**
  Laravel is served by ` + "`php artisan serve`" + `, a single-process development server
  that handles one request at a time with no OPcache. Django runs one gunicorn
  worker with ` + "`--reload`" + `. Neither is what production looks like. See
  ` + "`RESULTS-PROD.md`" + ` for the same measurement against real application servers.
- **Source and database are bind-mounted**, so on macOS every file read goes
  through Docker Desktop's filesystem layer.
{{ else }}- **Every backend gets four CPUs**, and the single-threaded runtimes are
  configured with four workers to match: a PHP-FPM pool of four, four gunicorn
  workers, four clustered Node processes, and ` + "`GOMAXPROCS=4`" + ` for Go. Without
  that cap Go would simply take all twelve cores of the host.
- **Nothing is bind-mounted.** The seeded SQLite database is copied into the
  image at build time, so no request touches the host filesystem.
- **This is still a single container per backend on a laptop.** It is not a
  tuned deployment, and it says nothing about how these stacks behave behind a
  load balancer, with a networked database, or under sustained traffic.
{{ end }}`

const comparisonTemplate = `# Velo-stats backends: development stack versus production stack

Generated {{ .GeneratedAt }} on {{ .Host }} with {{ .Docker }}.

Both stacks were measured in the same run, minutes apart, on the same machine.
The development profile is each repo's own ` + "`docker-compose.yml`" + `. The production
profile is the ` + "`docker-compose.prod.yml`" + ` alongside it: a real application server,
no bind mounts, the already-seeded SQLite database baked into the image, and a
four-CPU cap applied equally to every backend.

The point of this comparison is to separate two things that the development
numbers alone conflate: how fast a language and framework are, and how much of a
backend's measured slowness was really its development server.

| Setting | Value |
| --- | --- |
| Requests per endpoint per pass | {{ .Config.Requests }} |
| Concurrency in the concurrent pass | {{ .Config.Concurrency }} |
| Discarded warmup requests per endpoint | {{ .Config.Warmup }} |

## What changed per backend

| Backend | Development | Production |
| --- | --- | --- |
{{ range .Rows }}| {{ .Name }} | {{ .DevSetup }} | {{ .ProdSetup }} |
{{ end }}
## The gain, biggest first

| Backend | Dev latency | Prod latency | Latency change | Dev req/s | Prod req/s | Throughput change |
| --- | --- | --- | --- | --- | --- | --- |
{{ range .Rows }}| {{ .Name }} | {{ ms .DevLatency }} | {{ ms .ProdLatency }} | {{ speedup .LatencySpeedup }} | {{ rps .DevRPS }} | {{ rps .ProdRPS }} | {{ speedup .ThroughputSpeedup }} |
{{ end }}
## Throughput change per endpoint

Production requests per second divided by development requests per second, for
each endpoint.

| Endpoint |{{ range .Dev }} {{ .Backend.Name }} |{{ end }}
| --- |{{ range .Dev }} --- |{{ end }}
{{ $dev := .Dev }}{{ $prod := .Prod }}{{ range $endpoint := .Endpoints }}| ` + "`{{ $endpoint }}`" + ` |{{ range $i, $r := $dev }} {{ endpointSpeedup $dev $prod $i $endpoint }} |{{ end }}
{{ end }}
## Median latency side by side

| Endpoint |{{ range .Dev }} {{ .Backend.Name }} dev | {{ .Backend.Name }} prod |{{ end }}
| --- |{{ range .Dev }} --- | --- |{{ end }}
{{ range $endpoint := .Endpoints }}| ` + "`{{ $endpoint }}`" + ` |{{ range $i, $r := $dev }} {{ cell $r $endpoint }} | {{ if lt $i (len $prod) }}{{ cell (index $prod $i) $endpoint }}{{ else }}-{{ end }} |{{ end }}
{{ end }}
## Reading the gap

- **A large gain means the development stack was the bottleneck, not the
  language.** The clearest case is PHP: ` + "`php artisan serve`" + ` is a single-process
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
- **Full detail for each profile** is in ` + "`RESULTS-DEV.md`" + ` and
  ` + "`RESULTS-PROD.md`" + `, including per-endpoint percentiles and response sizes.
`
