package main

// Deployment is one way of running a backend: which Compose file describes it,
// which service exposes the HTTP port, and a one-line description that ends up
// in the report so the numbers are readable without opening the compose file.
type Deployment struct {
	ComposeFile string
	Service     string
	Description string
}

// Backend describes one implementation of the velo-stats API, in both of the
// shapes it can be run in.
type Backend struct {
	Name     string // short name, also the value accepted by -only
	Language string // shown in the report
	Dir      string // path to the repo checkout, relative to this harness
	BaseURL  string // where the harness reaches it from the host

	Dev  Deployment
	Prod Deployment

	// PathOverrides remaps a canonical endpoint path for this backend only.
	// Django mounts its list endpoints with a trailing slash and 301-redirects
	// the unslashed form, which would charge it for an extra round trip.
	PathOverrides map[string]string
}

// Deployment returns the deployment for the given profile.
func (b Backend) Deployment(profile string) Deployment {
	if profile == ProfileProd {
		return b.Prod
	}
	return b.Dev
}

// Path returns the URL path this backend serves the given canonical endpoint on.
func (b Backend) Path(endpoint string) string {
	if override, ok := b.PathOverrides[endpoint]; ok {
		return override
	}
	return endpoint
}

// The two profiles. Development is each repo's own docker-compose.yml, the
// stack a contributor runs. Production is docker-compose.prod.yml, added
// alongside it: real application servers, no bind mounts, the seeded database
// baked into the image, and the same CPU allowance for every backend.
const (
	ProfileDev  = "dev"
	ProfileProd = "prod"
)

// endpoints are the canonical paths, in the order they appear in the report.
// Every backend exposes all five as parameterless GETs.
var endpoints = []string{
	"/_healthcheck",
	"/rides",
	"/rides/summary",
	"/rides/cost",
	"/stations",
}

// repoRoot is the directory holding all the velo-stats checkouts. It is the
// parent of this harness's own directory.
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
