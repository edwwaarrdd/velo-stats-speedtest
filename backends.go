package main

type Deployment struct {
	ComposeFile string
	Service     string
	Description string
}

type Backend struct {
	Name     string // short name, also the value accepted by -only
	Language string // shown in the report
	Dir      string // path to the repo checkout, relative to this harness
	BaseURL  string // where the harness reaches it from the host

	Dev  Deployment
	Prod Deployment

	// Django mounts its list endpoints with a trailing slash and 301-redirects
	// the unslashed form, which would charge it for an extra round trip.
	PathOverrides map[string]string
}

func (b Backend) Deployment(profile string) Deployment {
	if profile == ProfileProd {
		return b.Prod
	}
	return b.Dev
}

func (b Backend) Path(endpoint string) string {
	if override, ok := b.PathOverrides[endpoint]; ok {
		return override
	}
	return endpoint
}

// Development is each repo's own docker-compose.yml. Production is
// docker-compose.prod.yml alongside it: real application servers, no bind
// mounts, the seeded database baked into the image, and the same CPU allowance
// for every backend.
const (
	ProfileDev  = "dev"
	ProfileProd = "prod"
)

// /_healthcheck is deliberately excluded: it does no work (see
// HealthcheckHandler and its equivalents in the other backends), so
// benchmarking it measures framework dispatch overhead rather than anything
// the application does. It is still used on its own by waitForHealthy to
// decide when a backend is ready.
var endpoints = []string{
	"/rides",
	"/rides/summary",
	"/rides/cost",
	"/stations",
}

const repoRoot = ".."

var backends = []Backend{
	{
		Name:     "golang",
		Language: "Go 1.25",
		Dir:      repoRoot + "/velo-stats-golang",
		BaseURL:  "http://127.0.0.1:8000",
		Dev: Deployment{
			ComposeFile: "docker-compose.yml",
			Service:     "app",
			Description: "net/http, bind-mounted SQLite, all cores",
		},
		Prod: Deployment{
			ComposeFile: "docker-compose.prod.yml",
			Service:     "app",
			Description: "net/http, baked-in SQLite, GOMAXPROCS=4",
		},
	},
	{
		Name:     "java",
		Language: "Spring Boot 4.1 / Java 25",
		Dir:      repoRoot + "/velo-stats-java",
		BaseURL:  "http://127.0.0.1:8000",
		Dev: Deployment{
			ComposeFile: "docker-compose.yml",
			Service:     "app",
			Description: "embedded server, bind-mounted SQLite, all cores",
		},
		Prod: Deployment{
			ComposeFile: "docker-compose.prod.yml",
			Service:     "app",
			Description: "embedded server, baked-in SQLite, ActiveProcessorCount=4",
		},
	},
	{
		Name:     "nodejs",
		Language: "NestJS 11 / Node 22",
		Dir:      repoRoot + "/velo-stats-nodejs",
		BaseURL:  "http://127.0.0.1:8000",
		Dev: Deployment{
			ComposeFile: "docker-compose.yml",
			Service:     "app",
			Description: "single Node process, bind-mounted source and SQLite",
		},
		Prod: Deployment{
			ComposeFile: "docker-compose.prod.yml",
			Service:     "app",
			Description: "4 clustered Node processes, baked-in SQLite",
		},
	},
	{
		Name:     "laravel",
		Language: "Laravel 13 / PHP 8.5",
		Dir:      repoRoot + "/velo-stats-laravel",
		BaseURL:  "http://127.0.0.1:8000",
		Dev: Deployment{
			ComposeFile: "docker-compose.yml",
			Service:     "app",
			Description: "php artisan serve, one request at a time, no OPcache",
		},
		Prod: Deployment{
			ComposeFile: "docker-compose.prod.yml",
			Service:     "web",
			Description: "nginx to a 4-worker PHP-FPM pool, OPcache with JIT",
		},
	},
	{
		Name:     "symfony",
		Language: "Symfony 8.1 / PHP 8.5",
		Dir:      repoRoot + "/velo-stats-symfony",
		BaseURL:  "http://127.0.0.1:8000",
		Dev: Deployment{
			ComposeFile: "docker-compose.yml",
			Service:     "app",
			Description: "php -S, one request at a time, no OPcache",
		},
		Prod: Deployment{
			ComposeFile: "docker-compose.prod.yml",
			Service:     "web",
			Description: "nginx to a 4-worker PHP-FPM pool, OPcache with JIT",
		},
	},
	{
		Name:     "python",
		Language: "Django 6 / gunicorn",
		Dir:      repoRoot + "/velo-stats-python",
		BaseURL:  "http://127.0.0.1:8000",
		Dev: Deployment{
			ComposeFile: "docker-compose.yml",
			Service:     "api",
			Description: "gunicorn, 1 worker, --reload, bind-mounted source",
		},
		Prod: Deployment{
			ComposeFile: "docker-compose.prod.yml",
			Service:     "api",
			Description: "gunicorn, 4 preloaded workers, baked-in SQLite",
		},
		PathOverrides: map[string]string{
			"/rides":    "/rides/",
			"/stations": "/stations/",
		},
	},
}
