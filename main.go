// Command speedtest benchmarks the four velo-stats backend implementations
// against each other and writes Markdown reports.
//
// It runs one backend at a time: bring the containers up, wait for the health
// check, warm it up, measure every endpoint sequentially and then concurrently,
// tear it down, move on. They cannot run together anyway, since all four publish
// port 8000, and running them alone is also the only way to compare them fairly.
//
// Each backend can be run in two shapes. The development profile is the
// docker-compose.yml a contributor uses. The production profile is the
// docker-compose.prod.yml alongside it: real application servers, no bind
// mounts, the seeded SQLite database baked into the image, and the same CPU
// allowance for every backend. Running both answers two questions at once -
// which backend is fastest, and how much of a backend's development number was
// the development stack rather than the language.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Config holds the run's tunables, all settable from the command line.
type Config struct {
	Requests    int
	Concurrency int
	Warmup      int
	Only        string
	Profile     string
	SkipBuild   bool
	OutDir      string
	BootTimeout time.Duration
}

func main() {
	var cfg Config
	flag.IntVar(&cfg.Requests, "requests", 100, "requests per endpoint, per pass")
	flag.IntVar(&cfg.Concurrency, "concurrency", 10, "parallel requests in the concurrent pass")
	flag.IntVar(&cfg.Warmup, "warmup", 20, "discarded warmup requests per endpoint")
	flag.StringVar(&cfg.Only, "only", "", "run a single backend by name (golang, nodejs, php, python)")
	flag.StringVar(&cfg.Profile, "profile", "both", "which stack to measure: dev, prod, or both")
	flag.BoolVar(&cfg.SkipBuild, "skip-build", false, "skip docker compose --build, reuse the existing image")
	flag.StringVar(&cfg.OutDir, "out", ".", "directory to write the Markdown reports into")
	flag.DurationVar(&cfg.BootTimeout, "boot-timeout", 5*time.Minute, "how long to wait for a backend to become healthy")
	flag.Parse()

	selected, err := selectBackends(cfg.Only)
	if err != nil {
		log.Fatal(err)
	}
	profiles, err := selectProfiles(cfg.Profile)
	if err != nil {
		log.Fatal(err)
	}

	// results[profile] holds one entry per backend, in backend order.
	results := map[string][]Result{}
	first := true

	for _, profile := range profiles {
		for _, b := range selected {
			if !first {
				// Let the previous container finish releasing port 8000.
				time.Sleep(3 * time.Second)
			}
			first = false
			results[profile] = append(results[profile], run(b, profile, cfg))
		}
	}

	var written []string
	for _, profile := range profiles {
		path := reportPath(cfg.OutDir, profile)
		if err := writeReport(path, cfg, profile, results[profile]); err != nil {
			log.Fatalf("writing %s: %v", path, err)
		}
		written = append(written, path)
	}

	// The comparison only means anything when both stacks were measured in the
	// same run, on the same machine, minutes apart.
	if len(profiles) == 2 {
		path := cfg.OutDir + "/RESULTS-DEV-VS-PROD.md"
		if err := writeComparison(path, cfg, results[ProfileDev], results[ProfileProd]); err != nil {
			log.Fatalf("writing %s: %v", path, err)
		}
		written = append(written, path)
	}

	logf("\nWrote %s", strings.Join(written, ", "))
}

// run benchmarks one backend in one profile from a cold start, and always tears
// it down again.
func run(b Backend, profile string, cfg Config) Result {
	d := b.Deployment(profile)
	logf("\n=== %s / %s (%s) ===", b.Name, profile, d.Description)

	logf("  starting containers...")
	if err := startBackend(b, d, !cfg.SkipBuild); err != nil {
		stopBackend(b, d)
		return failure(b, profile, d, fmt.Sprintf("could not start containers: %v", err))
	}
	defer func() {
		logf("  stopping containers...")
		if err := stopBackend(b, d); err != nil {
			logf("  warning: teardown failed: %v", err)
		}
	}()

	logf("  waiting for health check...")
	if err := waitForHealthy(b, cfg.BootTimeout); err != nil {
		return failure(b, profile, d, err.Error())
	}

	return benchmarkBackend(b, profile, cfg)
}

func failure(b Backend, profile string, d Deployment, note string) Result {
	logf("  FAILED: %s", note)
	return Result{Backend: b, Profile: profile, Deployment: d, Failed: true, Note: note}
}

// reportPath names the report for a profile.
func reportPath(dir, profile string) string {
	if profile == ProfileProd {
		return dir + "/RESULTS-PROD.md"
	}
	return dir + "/RESULTS-DEV.md"
}

// selectBackends returns every backend, or just the one named by -only.
func selectBackends(only string) ([]Backend, error) {
	if only == "" {
		return backends, nil
	}
	for _, b := range backends {
		if b.Name == only {
			return []Backend{b}, nil
		}
	}

	names := make([]string, len(backends))
	for i, b := range backends {
		names[i] = b.Name
	}
	return nil, fmt.Errorf("unknown backend %q, expected one of: %s", only, strings.Join(names, ", "))
}

// selectProfiles turns the -profile flag into the profiles to measure, in the
// order they are run.
func selectProfiles(profile string) ([]string, error) {
	switch profile {
	case "both":
		return []string{ProfileDev, ProfileProd}, nil
	case ProfileDev, ProfileProd:
		return []string{profile}, nil
	default:
		return nil, fmt.Errorf("unknown profile %q, expected dev, prod or both", profile)
	}
}

// logf writes progress to stderr, keeping stdout free of harness chatter.
func logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
