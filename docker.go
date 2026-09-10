package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// compose runs `docker compose -f <file> <args...>` inside the backend's repo
// directory and streams its output to our stderr, so a failing build is visible.
func compose(b Backend, d Deployment, args ...string) error {
	full := append([]string{"compose", "-f", d.ComposeFile}, args...)

	cmd := exec.Command("docker", full...)
	cmd.Dir = b.Dir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// startBackend brings up only the service that exposes the HTTP port. Compose
// pulls in whatever it depends on, so redis and, in the Laravel production
// stack, PHP-FPM come along. The three queue workers each repo defines are left
// down on purpose: they do no work during a benchmark, but they compete for CPU
// while they boot.
func startBackend(b Backend, d Deployment, build bool) error {
	args := []string{"up", "-d"}
	if build {
		args = append(args, "--build")
	}
	return compose(b, d, append(args, d.Service)...)
}

// stopBackend tears the whole project down, so the next backend starts from a
// clean host with port 8000 free.
func stopBackend(b Backend, d Deployment) error {
	return compose(b, d, "down", "--remove-orphans")
}

// waitForHealthy polls /_healthcheck until it answers 200 or the deadline passes.
func waitForHealthy(b Backend, timeout time.Duration) error {
	client := &http.Client{Timeout: 3 * time.Second}
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		resp, err := client.Get(b.BaseURL + b.Path("/_healthcheck"))
		if err != nil {
			lastErr = err
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("%s never became healthy within %s (last error: %v)", b.Name, timeout, lastErr)
}

// dockerVersion is recorded in the report so a stale result is recognisable.
func dockerVersion() string {
	out, err := exec.Command("docker", "--version").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}
