package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type EndpointResult struct {
	Endpoint string // canonical path, e.g. /rides
	URL      string // what was actually requested

	Status    int   // status of the single preflight request
	BodyBytes int64 // size of the preflight response body

	Sequential Stats
	Concurrent Stats
}

type Result struct {
	Backend    Backend
	Profile    string
	Deployment Deployment
	Endpoints  []EndpointResult

	Failed bool
	Note   string
}

// Keep-alives and idle connections are sized for the concurrent pass, so we
// measure the server rather than TCP handshakes. Redirects are refused: a 3xx
// means we requested the wrong path.
func newClient(concurrency int) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = concurrency * 2
	transport.MaxIdleConnsPerHost = concurrency * 2
	transport.DisableCompression = true

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// The body is drained, so the server's full response time is measured.
func doRequest(client *http.Client, url string) (time.Duration, int, int64, error) {
	start := time.Now()

	resp, err := client.Get(url)
	if err != nil {
		return time.Since(start), 0, 0, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(io.Discard, resp.Body)
	elapsed := time.Since(start)
	if err != nil {
		return elapsed, resp.StatusCode, n, err
	}
	if resp.StatusCode != http.StatusOK {
		return elapsed, resp.StatusCode, n, fmt.Errorf("status %d", resp.StatusCode)
	}
	return elapsed, resp.StatusCode, n, nil
}

func benchmarkBackend(b Backend, profile string, cfg Config) Result {
	client := newClient(cfg.Concurrency)
	result := Result{Backend: b, Profile: profile, Deployment: b.Deployment(profile)}

	for _, endpoint := range endpoints {
		url := b.BaseURL + b.Path(endpoint)
		er := EndpointResult{Endpoint: endpoint, URL: url}

		// Preflight: one request that records what this endpoint actually
		// serves, so the report can flag an unseeded database or a redirect.
		_, status, size, err := doRequest(client, url)
		er.Status = status
		er.BodyBytes = size
		if err != nil {
			logf("  %-16s preflight failed: %v", endpoint, err)
		}

		// Warmup, discarded: pays for JIT, opcode caching, connection setup and
		// the SQLite page cache before anything is measured.
		for i := 0; i < cfg.Warmup; i++ {
			_, _, _, _ = doRequest(client, url)
		}

		er.Sequential = runSequential(client, url, cfg.Requests)
		er.Concurrent = runConcurrent(client, url, cfg.Requests, cfg.Concurrency)

		logf("  %-16s seq median %-10s conc %7.1f req/s  (%d errors)",
			endpoint,
			formatDuration(er.Sequential.Median),
			er.Concurrent.RPS,
			er.Sequential.Errors+er.Concurrent.Errors)

		result.Endpoints = append(result.Endpoints, er)
	}
	return result
}

// One request after another: pure per-request latency with no queuing of our
// own making.
func runSequential(client *http.Client, url string, n int) Stats {
	durations := make([]time.Duration, 0, n)
	errors := 0

	start := time.Now()
	for i := 0; i < n; i++ {
		d, _, _, err := doRequest(client, url)
		if err != nil {
			errors++
			continue
		}
		durations = append(durations, d)
	}
	return summarise(durations, errors, time.Since(start))
}

// Spreading the requests over workers is where a single-process server
// separates from a multi-threaded one.
func runConcurrent(client *http.Client, url string, n, concurrency int) Stats {
	type outcome struct {
		duration time.Duration
		err      error
	}
	outcomes := make([]outcome, n)

	jobs := make(chan int)
	var wg sync.WaitGroup

	start := time.Now()
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				d, _, _, err := doRequest(client, url)
				outcomes[i] = outcome{d, err}
			}
		}()
	}
	for i := 0; i < n; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	elapsed := time.Since(start)

	durations := make([]time.Duration, 0, n)
	errors := 0
	for _, o := range outcomes {
		if o.err != nil {
			errors++
			continue
		}
		durations = append(durations, o.duration)
	}
	return summarise(durations, errors, elapsed)
}
