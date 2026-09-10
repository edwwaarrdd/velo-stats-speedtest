package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func compose(b Backend, d Deployment, args ...string) error {
	full := append([]string{"compose", "-f", d.ComposeFile}, args...)

	cmd := exec.Command("docker", full...)
	cmd.Dir = b.Dir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Only the service exposing the HTTP port is brought up; Compose pulls in what
// it depends on. The three queue workers each repo defines are left down on
// purpose: they do no work during a benchmark, but they compete for CPU while
// they boot.
func startBackend(b Backend, d Deployment, build bool) error {
	args := []string{"up", "-d"}
	if build {
		args = append(args, "--build")
	}
	return compose(b, d, append(args, d.Service)...)
}

// The whole project goes down, so the next backend starts with port 8000 free.
func stopBackend(b Backend, d Deployment) error {
	return compose(b, d, "down", "--remove-orphans")
}

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

// Recorded in the report, so a stale result is recognisable.
func dockerVersion() string {
	out, err := exec.Command("docker", "--version").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}
